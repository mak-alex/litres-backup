package litres

import (
	"encoding/json"
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

const (
	TB = 1000000000000
	GB = 1000000000
	MB = 1000000
	KB = 1000

	// urls
	baseUrl      = "http://robot.litres.ru/"
	authorizeUrl = baseUrl + "pages/catalit_authorise/"
	//genresUrl       = baseUrl + "pages/catalit_genres/"
	//authorsUrl      = baseUrl + "pages/catalit_persons/"
	catalogUrl = baseUrl + "pages/catalit_browser/"
	//trialsUrl       = baseUrl + "static/trials/"
	//purchaseUrl     = baseUrl + "pages/purchase_book/"
	downloadBookUrl = baseUrl + "pages/catalit_download_book/"
)

var (
	formats = []string{
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
		"fb3",
	}
)

type Litres struct {
	Login              string
	Password           string
	BookPath           string
	Format             string
	NormalizedName     bool
	Available4Download bool
	Progress           bool
	Verbose            bool
	Debug              bool
	sid                string
}

func New(litres *Litres) *Litres {
	if !litres.Available4Download {
		if strings.EqualFold(litres.Format, "list") || !litres.existFormat() {
			litres.printFormats()
		}
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

func (l *Litres) ShowAvailable4Download(books []model.Fb2Book) {
	fmt.Println("Display a list of available books for download:")
	for i, book := range books {
		filename := fmt.Sprintf(
			"\t%d. %s %s - %s",
			i,
			book.TextDescription.Hidden.TitleInfo.Author.FirstName,
			book.TextDescription.Hidden.TitleInfo.Author.LastName,
			book.TextDescription.Hidden.TitleInfo.BookTitle,
		)
		fmt.Println(filename)
	}
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
	if r == nil {
		return
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil && l.Verbose {
		log.Fatal(err)
	}
	if res == nil {
		return
	}
	if l.Debug {
		log.Println(res.Status)
	}
	defer func() {
		_ = res.Body.Close()
	}()
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

func (l *Litres) GetBooks(checkpoint, search *string) *model.CatalitFb2Books {
	data := url.Values{}
	data.Set("sid", l.sid)
	data.Set("my", "1")
	if checkpoint != nil && *checkpoint != "" {
		data.Set("checkpoint", *checkpoint)
	}
	if search != nil && *search != "" {
		data.Set("search", *search)
	}
	data.Set("limit", "0,1000")

	client := &http.Client{}
	r, err := http.NewRequest("POST", catalogUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil && l.Verbose {
		log.Fatal(err)
	}
	if r == nil {
		return nil
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil && l.Verbose {
		log.Fatal(err)
	}
	if res == nil {
		return nil
	}
	if l.Debug {
		log.Println(res.Status)
	}
	defer func() {
		_ = res.Body.Close()
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil && l.Verbose {
		log.Fatal(err)
	}
	catalitFb2Books := model.CatalitFb2Books{}
	if err := xml.Unmarshal(body, &catalitFb2Books); err != nil && l.Verbose {
		log.Fatal(err)
	}
	if l.Debug {
		l.prettyPrint(catalitFb2Books)
	}

	return &catalitFb2Books
}

func (l *Litres) DownloadBooks(checkpoint, search *string, id *int) ([]string, error) {
	list := l.GetBooks(checkpoint, search)

	if l.Available4Download {
		l.ShowAvailable4Download(list.Fb2Book)
		os.Exit(0)
	}

	done := make(chan string, len(list.Fb2Book))
	errch := make(chan error, len(list.Fb2Book))

	if id == nil || *id == -1 {
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

		if err == nil {
			log.Println("Total downloaded books:", list.Records)
		}
		return bytesArray, err
	} else {
		file := list.Fb2Book[*id]
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

		_, err := l.download(file.HubID, path.Join(l.BookPath, filename))
		return nil, err
	}
}

func (l *Litres) existsBook(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (l *Litres) getFileSize1(filePath string) (int64, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	// get the size
	return fi.Size(), nil
}

func (l *Litres) getFileSize2(filePath string) (int64, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = f.Close()
	}()
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func (l *Litres) download(hubID, filePath string) (body string, err error) {

	data := url.Values{}
	data.Set("sid", l.sid)
	data.Set("art", hubID)
	data.Set("type", l.Format)

	client := &http.Client{}

	r, err := http.NewRequest("POST", downloadBookUrl, strings.NewReader(data.Encode()))
	if err != nil && l.Verbose {
		log.Fatal(err)
	}
	if r == nil {
		return
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil && l.Verbose {
		log.Fatal(err)
	}
	if res == nil {
		return
	}
	if l.Debug {
		log.Println(res.Status)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if !strings.Contains(res.Header.Get("Content-Disposition"), "attachment") {
		log.Printf("Downloading this book: %s is not possible", filePath)
		return "", errors.New("<catalit-download-drm-failed/>")
	}

	fsize, _ := strconv.Atoi(res.Header.Get("Content-Length"))

	if localFileSize, err := l.getFileSize1(filePath); err == nil {
		if l.existsBook(filePath) && localFileSize == int64(fsize) {
			log.Printf("Book %s (%s) exists", filePath, l.lenReadable(fsize, 2))
			return "", err
		}
	}

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = out.Close()
	}()

	if l.Progress {
		counter := bar.NewWriteCounter(fsize, filePath)
		counter.Start()
		_, err = io.Copy(out, io.TeeReader(res.Body, counter))
		if err != nil {
			return "", err
		}

		counter.Finish()
	} else {
		log.Printf("Book: %s downloaded", filePath)
	}

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return "", err
	}

	return filepath.Base(strings.TrimSpace(out.Name())), nil
}

func (l *Litres) lenReadable(length int, decimals int) (out string) {
	var unit string
	var i int
	var remainder int

	// Get whole number, and the remainder for decimals
	if length > TB {
		unit = "TB"
		i = length / TB
		remainder = length - (i * TB)
	} else if length > GB {
		unit = "GB"
		i = length / GB
		remainder = length - (i * GB)
	} else if length > MB {
		unit = "MB"
		i = length / MB
		remainder = length - (i * MB)
	} else if length > KB {
		unit = "KB"
		i = length / KB
		remainder = length - (i * KB)
	} else {
		return strconv.Itoa(length) + " B"
	}

	if decimals == 0 {
		return strconv.Itoa(i) + " " + unit
	}

	// This is to calculate missing leading zeroes
	width := 0
	if remainder > GB {
		width = 12
	} else if remainder > MB {
		width = 9
	} else if remainder > KB {
		width = 6
	} else {
		width = 3
	}

	// Insert missing leading zeroes
	remainderString := strconv.Itoa(remainder)
	for iter := len(remainderString); iter < width; iter++ {
		remainderString = "0" + remainderString
	}
	if decimals > len(remainderString) {
		decimals = len(remainderString)
	}

	return fmt.Sprintf("%d.%s %s", i, remainderString[:decimals], unit)
}

func (l *Litres) prettyPrint(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(b))
}
