package litres

import (
	"bytes"
	"fmt"
	"github.com/mak-alex/litres-backup/pkg/bar"
	"github.com/mak-alex/litres-backup/pkg/consts"
	"github.com/mak-alex/litres-backup/pkg/logger"
	"github.com/mak-alex/litres-backup/tools"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Litres struct {
	Login              string
	Password           string
	Library            string
	Format             string
	NormalizedName     bool
	Available4Download bool
	Progress           bool
	Verbose            bool
	Debug              bool
	sid                string
	tmpFile            string
}

func New(litres *Litres) *Litres {
	if !litres.Available4Download {
		if strings.EqualFold(litres.Format, "list") || !litres.existFormat() {
			litres.printFormats()
		}
	}

	if litres.Login == "" {
		logger.Work.Error("[litres.New] can't be nil", zap.String("login", litres.Login))
	}
	if litres.Password == "" {
		logger.Work.Error("[litres.New] can't be nil", zap.String("password", litres.Password))
	}
	if litres.Library == "" {
		logger.Work.Error("[litres.New] can't be nil", zap.String("library", litres.Library))
	}

	litres.tmpFile = filepath.Join(litres.Library, "litres.xml")

	return litres
}

func (l *Litres) download(hubID, filePath string, fileSize int) (bodyBook string, err error) {
	l.authorization()

	data := url.Values{}
	data.Set("sid", l.sid)
	data.Set("art", hubID)
	data.Set("type", l.Format)

	client := &http.Client{}

	r, err := http.NewRequest("POST", consts.DownloadBookUrl, strings.NewReader(data.Encode()))
	if err != nil {
		logger.Work.Fatal("[litres.download]", zap.Error(err))
	}
	if r == nil {
		return
	}
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36 Edg/89.0.774.68")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		logger.Work.Fatal("[litres.download]", zap.Error(err))
	}
	if res == nil {
		return
	}
	defer func() {
		_ = res.Body.Close()
	}()

	body, _ := ioutil.ReadAll(res.Body)
	if strings.Contains(string(body), "catalit-download-drm-failed") {
		logger.Work.Error("At the request of the copyright holder, this book is not available for download as a file", zap.String("name", filepath.Base(filePath)), zap.String("hubID", hubID))
		tools.OpenBrowser(fmt.Sprintf("http://www.litres.ru/pages/biblio_book/?art=%s", hubID))
		return
	}

	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	//fsize, _ := strconv.Atoi(res.Header.Get("Content-Length"))
	if localFileSize, err := tools.GetFileSize(filePath); err == nil {
		if !tools.FileNotExists(filePath) && localFileSize == int64(fileSize) {
			logger.Work.Info("[litres.download] exists", zap.String("filePath", filePath), zap.String("size", tools.LenReadable(fileSize, 2)))
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
		counter := bar.NewWriteCounter(fileSize, filePath)
		counter.Start()
		_, err = io.Copy(out, io.TeeReader(res.Body, counter))
		if err != nil {
			return "", err
		}

		counter.Finish()
	} else {
		logger.Work.Info("[litres.download] downloaded", zap.String("filePath", filePath))
	}

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return "", err
	}

	return filepath.Base(strings.TrimSpace(out.Name())), nil
}
