package litres

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mak-alex/litres-backup/pkg/consts"
	"github.com/mak-alex/litres-backup/pkg/logger"
	"github.com/mak-alex/litres-backup/pkg/model"
	"github.com/mak-alex/litres-backup/tools"
	"go.uber.org/zap"
)

func (l *Litres) LoadPurchasedBooks(search *string, offset, maxCount int) (result model.CatalitFb2Books) {
	params := url.Values{}
	params.Set("my", "1")
	params.Set("limit", fmt.Sprintf("%d,%d", offset, maxCount))
	if search != nil && *search != "" {
		params.Set("search", *search)
	}

	l.loadBooks(params, &result)
	return
}

func (l *Litres) LoadPopularBooks(offset, maxCount int) (result model.CatalitFb2Books) {
	params := url.Values{}
	params.Set("rating", "books")
	params.Set("limit", fmt.Sprintf("%d,%d", offset, maxCount))

	l.loadBooks(params, &result)
	return
}

func (l *Litres) LoadNewBooks(offset, maxCount int) (result model.CatalitFb2Books) {
	params := url.Values{}
	params.Set("rating", "hot")
	params.Set("limit", fmt.Sprintf("%d,%d", offset, maxCount))

	l.loadBooks(params, &result)
	return
}

func (l *Litres) LoadBooksByGenre(genreID, offset, maxCount int) (result model.CatalitFb2Books) {
	params := url.Values{}
	params.Set("genre", strconv.Itoa(genreID))
	params.Set("limit", fmt.Sprintf("%d,%d", offset, maxCount))

	l.loadBooks(params, &result)
	return
}

func (l *Litres) LoadBooksByBookId(bookID int, myOnly bool) (result model.CatalitFb2Books) {
	params := url.Values{}
	params.Set("art", strconv.Itoa(bookID))
	if myOnly {
		params.Set("my", "1")
	}

	l.loadBooks(params, &result)
	return
}

func (l *Litres) LoadBooksByAuthor(authorID, offset, maxCount int) (result model.CatalitFb2Books) {
	params := url.Values{}

	params.Set("person", strconv.Itoa(authorID))

	l.loadBooks(params, &result)
	return
}

func (l *Litres) loadBooks(params url.Values, result *model.CatalitFb2Books) (err error) {
	var body []byte

	if tools.FileNotExists(l.tmpFile) {
		l.authorization()
		params.Set("sid", l.sid)
		params.Set("search_types", "0")
		params.Set("checkpoint", "2000-01-01 00:00:00")

		client := &http.Client{}
		r, err := http.NewRequest("POST", consts.CatalogUrl, strings.NewReader(params.Encode()))
		if err != nil && l.Verbose {
			logger.Work.Fatal("[litres.GetBooks]", zap.Error(err))
		}
		if r == nil {
			return nil
		}
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))

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
		if err != nil {
			logger.Work.Fatal("[litres.GetBooks]", zap.Error(err))
		}
	}

	err = xml.Unmarshal(body, result)
	if err != nil && l.Verbose {
		logger.Work.Fatal("[litres.GetBooks]", zap.Error(err))
	}

	if l.Debug {
		tools.PrettyPrint(result)
	}

	if !l.Available4Download {
		l.downloadBooks(result, nil, nil)
	}
	return
}

func (l *Litres) downloadBooks(list *model.CatalitFb2Books, title, id *string) ([]string, error) {
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
