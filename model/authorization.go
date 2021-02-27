package model

import "encoding/xml"

type CatalitAuthorizationOk struct {
	XMLName      xml.Name `xml:"catalit-authorization-ok"`
	Text         string   `xml:",chardata"`
	UserID       string   `xml:"user-id,attr"`
	FirstName    string   `xml:"first-name,attr"`
	LastName     string   `xml:"last-name,attr"`
	Login        string   `xml:"login,attr"`
	Mail         string   `xml:"mail,attr"`
	Phone        string   `xml:"phone,attr"`
	Tcountry     string   `xml:"tcountry,attr"`
	CanRebill    string   `xml:"can-rebill,attr"`
	BirthDay     string   `xml:"birth_day,attr"`
	Male         string   `xml:"male,attr"`
	Now          string   `xml:"now,attr"`
	Sid          string   `xml:"sid,attr"`
	BooksCnt     string   `xml:"books-cnt,attr"`
	AuthorsCnt   string   `xml:"authors-cnt,attr"`
	BiblioUser   string   `xml:"biblio_user,attr"`
	Account      string   `xml:"account,attr"`
	Bonus        string   `xml:"bonus,attr"`
	AccountFull  string   `xml:"account_full,attr"`
	SocialLogins struct {
		Text        string `xml:",chardata"`
		SocialLogin []struct {
			Text    string `xml:",chardata"`
			Service string `xml:"service,attr"`
			UserID  string `xml:"user-id,attr"`
			Name    string `xml:"name,attr"`
		} `xml:"social-login"`
	} `xml:"social-logins"`
}
