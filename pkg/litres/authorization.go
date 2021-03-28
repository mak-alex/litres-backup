package litres

import (
	"encoding/xml"
	"github.com/mak-alex/litres-backup/pkg/consts"
	"github.com/mak-alex/litres-backup/pkg/logger"
	"github.com/mak-alex/litres-backup/pkg/model"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		logger.Work.Fatal("[litres.authorization]", zap.Error(err))
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
	if err != nil {
		logger.Work.Fatal("[litres.authorization]", zap.Error(err))
	}
	if l.Debug {
		logger.Work.Debug("[litres.authorization]", zap.Any("response", body))
	}
	catalitAuthorizationOk := model.CatalitAuthorizationOk{}
	if err := xml.Unmarshal(body, &catalitAuthorizationOk); err != nil {
		logger.Work.Fatal("[litres.authorization]", zap.Error(err))
	}

	l.sid = catalitAuthorizationOk.Sid
}
