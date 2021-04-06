package litres

import (
	"errors"
	"github.com/mak-alex/litres-backup/pkg/bar"
	"github.com/mak-alex/litres-backup/pkg/consts"
	"github.com/mak-alex/litres-backup/pkg/logger"
	"github.com/mak-alex/litres-backup/tools"
	"go.uber.org/zap"
	"io"
	"log"
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

func (l *Litres) download(hubID, filePath string) (body string, err error) {
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
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		logger.Work.Fatal("[litres.download]", zap.Error(err))
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
		logger.Work.Info("[litres.download] not possible", zap.String("filePath", filePath))
		return "", errors.New("<catalit-download-drm-failed/>")
	}

	fsize, _ := strconv.Atoi(res.Header.Get("Content-Length"))

	if localFileSize, err := tools.GetFileSize(filePath); err == nil {
		if !tools.FileNotExists(filePath) && localFileSize == int64(fsize) {
			logger.Work.Info("[litres.download] exists", zap.String("filePath", filePath), zap.String("size", tools.LenReadable(fsize, 2)))
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
		logger.Work.Info("[litres.download] downloaded", zap.String("filePath", filePath))
	}

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return "", err
	}

	return filepath.Base(strings.TrimSpace(out.Name())), nil
}
