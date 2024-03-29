package litres

import (
	"encoding/xml"
	"fmt"
	"github.com/mak-alex/litres-backup/pkg/consts"
	"github.com/mak-alex/litres-backup/pkg/logger"
	"github.com/mak-alex/litres-backup/pkg/model"
	"github.com/mak-alex/litres-backup/tools"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func (l *Litres) GetBooks(checkpoint, search *string) *model.CatalitFb2Books {
	var (
		err  error
		body []byte
	)
	catalitFb2Books := model.CatalitFb2Books{}
	if tools.FileNotExists(l.tmpFile) {
		l.authorization()
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
		r, err := http.NewRequest("POST", consts.CatalogUrl, strings.NewReader(data.Encode())) // URL-encoded payload
		if err != nil && l.Verbose {
			logger.Work.Fatal("[litres.GetBooks]", zap.Error(err))
		}
		if r == nil {
			return nil
		}
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

		res, err := client.Do(r)
		if err != nil && l.Verbose {
			logger.Work.Fatal("[litres.GetBooks]", zap.Error(err))
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
		body, err = ioutil.ReadAll(res.Body)
		if err != nil && l.Verbose {
			logger.Work.Fatal("[litres.GetBooks]", zap.Error(err))
		}
		err = tools.WriteToFile(l.tmpFile, string(body))
		if err != nil {
			logger.Work.Fatal("[litres.GetBooks]", zap.Error(err))
		}
	} else {
		body, err = tools.ReadFile(l.tmpFile)
	}

	err = xml.Unmarshal(body, &catalitFb2Books)
	if err != nil && l.Verbose {
		logger.Work.Fatal("[litres.GetBooks]", zap.Error(err))
	}
	if l.Debug {
		tools.PrettyPrint(catalitFb2Books)
	}

	return &catalitFb2Books
}

func (l *Litres) DownloadBooks(checkpoint, title, id *string) ([]string, error) {
	list := l.GetBooks(checkpoint, title)
	do := func(file *model.Fb2Book) (string, error) {
		if !file.CheckFileType(l.Format) {
			return "", fmt.Errorf("it is not possible to download the book using this format %s", l.Format)
		}
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

		return l.download(file.HubID, path.Join(l.Library, filename), file.GetSizeByFileType(l.Format))
	}

	if id != nil && *id != "" {
		res, err := do(list.GetBookByID(id, &l.Format, true))
		return []string{res}, err
	} else if title != nil && *title != "" {
		res, err := do(list.GetBookByTitle(title, &l.Format, true))
		return []string{res}, err
	}

	books := make([]string, 0)
	for _, file := range list.Fb2Book {
		name, err := do(&file)
		if err != nil {
			return books, err
		}
		books = append(books, name)
	}

	return books, nil
}
