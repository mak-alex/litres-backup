package litres

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mak-alex/backlitr/bar"
	"github.com/mak-alex/backlitr/model"
)

var (
	baseUrl         = "http://robot.litres.ru/"
	authorizeUrl    = baseUrl + "pages/catalit_authorise/"
	genresUrl       = baseUrl + "pages/catalit_genres/"
	authorsUrl      = baseUrl + "pages/catalit_persons/"
	catalogUrl      = baseUrl + "pages/catalit_browser/"
	trialsUrl       = baseUrl + "static/trials/"
	purchaseUrl     = baseUrl + "pages/purchase_book/"
	downloadBookUrl = baseUrl + "pages/catalit_download_book/"
	formats         = []string{
		"fb2.zip",
		"html",
		"html.zip",
		"txt",
		"txt.zip",
		"rtf.zip",
		"a4.pdf",
		"a6.pdf",
		"mobi.prc",
		"epub",
		"ios.epub",
	}
)

type Litres struct {
	Login          string
	Password       string
	BookPath       string
	Format         string
	NormalizedName bool
	Progress       bool
	Verbose        bool
	Debug          bool
	sid            string
}

func New(litres *Litres) *Litres {
	if strings.EqualFold(litres.Format, "list") || !litres.existFormat() {
		litres.printFormats()
	}
	if litres.Login == "" {
		log.Fatal("'Login` can't be nil")
	}
	if litres.Password == "" {
		log.Fatal("'Password` can't be nil")
	}
	if litres.BookPath == "" {
		log.Fatal("'BookPath` can't be nil")
	}

	litres.authorise()

	return litres
}

func (l *Litres) existFormat() bool {
	for _, format := range formats {
		if strings.Contains(l.Format, format) {
			return true
		}
	}
	return false
}

func (l *Litres) printFormats() {
	fmt.Println("Available formats:")
	for _, format := range formats {
		fmt.Println("\t -", format)
	}
	os.Exit(0)
}

func (l *Litres) authorise() {
	data := url.Values{}
	data.Set("login", l.Login)
	data.Set("pwd", l.Password)

	if l.Debug {
		log.Printf("Authorization data: %v", data)
	}
	client := &http.Client{}
	r, err := http.NewRequest("POST", authorizeUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil && l.Verbose {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil && l.Verbose {
		log.Fatal(err)
	}
	if l.Debug {
		log.Println(res.Status)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil && l.Verbose {
		log.Fatal(err)
	}
	if l.Debug {
		log.Println(string(body))
	}
	catalitAuthorizationOk := model.CatalitAuthorizationOk{}
	if err := xml.Unmarshal(body, &catalitAuthorizationOk); err != nil {
		if l.Verbose {
			log.Fatal(err)
		} else {
			log.Fatal("Authorization failed")
		}
	}

	l.sid = catalitAuthorizationOk.Sid
}

func (l *Litres) GetBooks(checkpoint, search string) *model.CatalitFb2Books {
	data := url.Values{}
	data.Set("sid", l.sid)
	data.Set("my", "1")
	if checkpoint != "" {
		data.Set("checkpoint", checkpoint)
	}
	if search != "" {
		data.Set("search", search)
	}
	data.Set("limit", "0,1000")

	client := &http.Client{}
	r, err := http.NewRequest("POST", catalogUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil && l.Verbose {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil && l.Verbose {
		log.Fatal(err)
	}
	if l.Debug {
		log.Println(res.Status)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil && l.Verbose {
		log.Fatal(err)
	}
	if l.Debug {
		log.Println(string(body))
	}
	catalitFb2Books := model.CatalitFb2Books{}
	if err := xml.Unmarshal(body, &catalitFb2Books); err != nil && l.Verbose {
		log.Fatal(err)
	}

	log.Println("Total books:", catalitFb2Books.Records)

	return &catalitFb2Books
}

func (l *Litres) DownloadBooks(checkpoint, search string) ([]string, error) {
	list := l.GetBooks(checkpoint, search)

	done := make(chan string, len(list.Fb2Book))
	errch := make(chan error, len(list.Fb2Book))
	for _, file := range list.Fb2Book {
		go func(file model.Fb2Book) {
			ext := filepath.Ext(file.Filename)

			var filename string
			if l.NormalizedName {
				filename = fmt.Sprintf(
					"%s %s - %s.%s",
					file.TextDescription.Hidden.TitleInfo.Author.FirstName,
					file.TextDescription.Hidden.TitleInfo.Author.LastName,
					file.TextDescription.Hidden.TitleInfo.BookTitle,
					l.Format,
				)
			} else {
				filename = strings.ReplaceAll(file.Filename, strings.TrimLeft(ext, "."), l.Format)
			}
			if l.Debug {
				log.Println("Filename:", filename)
			}

			name, err := l.download(file.HubID, path.Join(l.BookPath, filename))
			if err != nil {
				errch <- err
				done <- ""
				return
			}
			done <- name
			errch <- nil
		}(file)
	}
	bytesArray := make([]string, 0)
	var errStr string
	for i := 0; i < len(list.Fb2Book); i++ {
		bytesArray = append(bytesArray, <-done)
		if err := <-errch; err != nil {
			errStr = errStr + " " + err.Error()
		}
	}
	var err error
	if errStr != "" {
		err = errors.New(errStr)
	}
	return bytesArray, err
}

func (l *Litres) existsBook(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (l *Litres) getFileSize1(filepath string) (int64, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	// get the size
	return fi.Size(), nil
}

func (l *Litres) getFileSize2(filepath string) (int64, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func (l *Litres) download(hubID, filepath string) (body string, err error) {

	data := url.Values{}
	data.Set("sid", l.sid)
	data.Set("art", hubID)
	data.Set("type", l.Format)

	client := &http.Client{}

	r, err := http.NewRequest("POST", downloadBookUrl, strings.NewReader(data.Encode()))
	if err != nil && l.Verbose {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil && l.Verbose {
		log.Fatal(err)
	}
	if l.Debug {
		log.Println(res.Status)
	}
	defer res.Body.Close()

	fsize, _ := strconv.Atoi(res.Header.Get("Content-Length"))

	localFileSize, err := l.getFileSize1(filepath)
	if err != nil {
		return
	}

	if l.existsBook(filepath) && localFileSize == int64(fsize) {
		return
	}

	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if l.Progress {
		counter := bar.NewWriteCounter(fsize, filepath)
		counter.Start()
		_, err = io.Copy(out, io.TeeReader(res.Body, counter))
		if err != nil {
			return "", err
		}

		counter.Finish()
	}

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return "", err
	}

	return out.Name(), nil
}
