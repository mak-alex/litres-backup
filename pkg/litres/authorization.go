package litres

import (
	"encoding/xml"
	"github.com/antchfx/htmlquery"
	"github.com/mak-alex/litres-backup/pkg/consts"
	"github.com/mak-alex/litres-backup/pkg/logger"
	"github.com/mak-alex/litres-backup/pkg/model"
	"github.com/mak-alex/litres-backup/tools"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func (l *Litres) authorization() {
	data := url.Values{}
	data.Set("login", l.Login)
	data.Set("pwd", l.Password)

	if l.Debug {
		logger.Work.Debug("[litres.authorization]", zap.Any("params", data))
	}
	client := &http.Client{}
	r, err := http.NewRequest("POST", consts.AuthorizeUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil && l.Verbose {
		logger.Work.Fatal("[litres.authorization]", zap.Error(err))
	}
	if r == nil {
		return
	}
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36 Edg/89.0.774.68")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		logger.Work.Fatal("[litres.authorization]", zap.Error(err))
	}

	if res == nil {
		return
	}

	defer func() {
		_ = res.Body.Close()
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Work.Fatal("[litres.authorization]", zap.Error(err))
	}

	if strings.Contains(res.Header.Get("Content-Type"), "text/html") {
		doc, err := htmlquery.Parse(strings.NewReader(string(body)))
		if err != nil {
			return
		}
		for _, attr := range htmlquery.Find(doc, "//form")[0].Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, "captcha_block_form") {
				logger.Work.Error("[litres.authorization] the lock worked, I can't bypass the captcha yet, attempt to open a page with a captcha in your browser", zap.String("link", consts.BaseUrl))
				tools.OpenBrowser(consts.BaseUrl)
				err := os.Remove(l.tmpFile)
				if err != nil {
					logger.Work.Error("[litres.remove.tmp] couldn't delete temporary file", zap.Error(err))
					return
				}
				logger.Work.Info("[litres.remove.tmp] temporary file", zap.String("deleted", l.tmpFile))
				return
			} else {
				logger.Work.Debug("[litres.authorization]", zap.Any("response", body))
			}
		}
	}

	catalitAuthorizationOk := model.CatalitAuthorizationOk{}
	if err := xml.Unmarshal(body, &catalitAuthorizationOk); err != nil {
		logger.Work.Fatal("[litres.authorization]", zap.Error(err))
		return
	}

	l.sid = catalitAuthorizationOk.Sid
}
