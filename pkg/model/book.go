package model

import "encoding/xml"

type CatalitFb2Books struct {
	XMLName     xml.Name  `xml:"catalit-fb2-books"`
	Text        string    `xml:",chardata"`
	Now         string    `xml:"now,attr"`
	Pages       string    `xml:"pages,attr"`
	Records     string    `xml:"records,attr"`
	Account     string    `xml:"account,attr"`
	AccountFull string    `xml:"account_full,attr"`
	Bonus       string    `xml:"bonus,attr"`
	OffersCnt   string    `xml:"offers_cnt,attr"`
	Fb2Book     []Fb2Book `xml:"fb2-book"`
}

type Fb2Book struct {
	Text              string `xml:",chardata"`
	HubID             string `xml:"hub_id,attr"`
	Added             string `xml:"added,attr"`
	Updated           string `xml:"updated,attr"`
	Chars             string `xml:"chars,attr"`
	Images            string `xml:"images,attr"`
	ZipSize           string `xml:"zip_size,attr"`
	HasTrial          string `xml:"has_trial,attr"`
	Type              string `xml:"type,attr"`
	Drm               string `xml:"drm,attr"`
	CoverH            string `xml:"cover_h,attr"`
	CoverW            string `xml:"cover_w,attr"`
	Filename          string `xml:"filename,attr"`
	Adult             string `xml:"adult,attr"`
	TrialPercent      string `xml:"trial_percent,attr"`
	Basket            string `xml:"basket,attr"`
	Payed             string `xml:"payed,attr"`
	Rating            string `xml:"rating,attr"`
	Recenses          string `xml:"recenses,attr"`
	BasePrice         string `xml:"base_price,attr"`
	Price             string `xml:"price,attr"`
	ItemsLeft         string `xml:"items_left,attr"`
	UserCloud         string `xml:"user_cloud,attr"`
	Voted1            string `xml:"voted1,attr"`
	Voted2            string `xml:"voted2,attr"`
	Voted3            string `xml:"voted3,attr"`
	Voted4            string `xml:"voted4,attr"`
	Voted5            string `xml:"voted5,attr"`
	UserVote          string `xml:"user_vote,attr"`
	Cover             string `xml:"cover,attr"`
	CoverPreview      string `xml:"cover_preview,attr"`
	URL               string `xml:"url,attr"`
	Copyright         string `xml:"copyright,attr"`
	ReadedPercent     string `xml:"readed_percent,attr"`
	BiblioSelfservice string `xml:"biblio_selfservice,attr"`
	BiblioFund        string `xml:"biblio_fund,attr"`
	BiblioBusy        string `xml:"biblio_busy,attr"`
	BiblioQueueSize   string `xml:"biblio_queue_size,attr"`
	BiblioMyRequest   string `xml:"biblio_my_request,attr"`
	InFolder          string `xml:"in_folder,attr"`
	TextDescription   struct {
		Text   string `xml:",chardata"`
		Hidden struct {
			Text      string `xml:",chardata"`
			TitleInfo struct {
				Text   string   `xml:",chardata"`
				Genre  []string `xml:"genre"`
				Author struct {
					Text       string `xml:",chardata"`
					FirstName  string `xml:"first-name"`
					LastName   string `xml:"last-name"`
					ID         string `xml:"id"`
					MiddleName string `xml:"middle-name"`
				} `xml:"author"`
				BookTitle  string `xml:"book-title"`
				Annotation struct {
					Text string `xml:",chardata"`
					P    []struct {
						Text     string `xml:",chardata"`
						Emphasis string `xml:"emphasis"`
					} `xml:"p"`
				} `xml:"annotation"`
				Keywords string `xml:"keywords"`
				Date     struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value,attr"`
				} `xml:"date"`
				Coverpage struct {
					Text  string `xml:",chardata"`
					Image struct {
						Text string `xml:",chardata"`
						Href string `xml:"href,attr"`
					} `xml:"image"`
				} `xml:"coverpage"`
				Lang       string `xml:"lang"`
				SrcLang    string `xml:"src-lang"`
				Translator []struct {
					Text       string `xml:",chardata"`
					FirstName  string `xml:"first-name"`
					LastName   string `xml:"last-name"`
					ID         string `xml:"id"`
					MiddleName string `xml:"middle-name"`
				} `xml:"translator"`
				Sequence []struct {
					Text     string `xml:",chardata"`
					Name     string `xml:"name,attr"`
					Number   string `xml:"number,attr"`
					Sequence struct {
						Text   string `xml:",chardata"`
						Name   string `xml:"name,attr"`
						Number string `xml:"number,attr"`
					} `xml:"sequence"`
				} `xml:"sequence"`
			} `xml:"title-info"`
			DocumentInfo struct {
				Text   string `xml:",chardata"`
				Author struct {
					Text      string `xml:",chardata"`
					Nickname  string `xml:"nickname"`
					FirstName string `xml:"first-name"`
					LastName  string `xml:"last-name"`
					Email     string `xml:"email"`
				} `xml:"author"`
				ProgramUsed string `xml:"program-used"`
				Date        struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value,attr"`
				} `xml:"date"`
				SrcURL  []string `xml:"src-url"`
				SrcOcr  string   `xml:"src-ocr"`
				ID      string   `xml:"id"`
				Version string   `xml:"version"`
				History struct {
					Text string   `xml:",chardata"`
					P    []string `xml:"p"`
				} `xml:"history"`
				Publisher []struct {
					Text      string `xml:",chardata"`
					FirstName string `xml:"first-name"`
					LastName  string `xml:"last-name"`
					ID        string `xml:"id"`
				} `xml:"publisher"`
			} `xml:"document-info"`
			PublishInfo struct {
				Text      string `xml:",chardata"`
				BookName  string `xml:"book-name"`
				Publisher string `xml:"publisher"`
				City      string `xml:"city"`
				Year      string `xml:"year"`
				ISBN      string `xml:"isbn"`
				Sequence  struct {
					Text   string `xml:",chardata"`
					Name   string `xml:"name,attr"`
					Number string `xml:"number,attr"`
				} `xml:"sequence"`
			} `xml:"publish-info"`
			SrcTitleInfo struct {
				Text   string `xml:",chardata"`
				Genre  string `xml:"genre"`
				Author struct {
					Text      string `xml:",chardata"`
					FirstName string `xml:"first-name"`
					LastName  string `xml:"last-name"`
				} `xml:"author"`
				BookTitle string `xml:"book-title"`
				Lang      string `xml:"lang"`
			} `xml:"src-title-info"`
			CustomInfo []struct {
				Text     string `xml:",chardata"`
				InfoType string `xml:"info-type,attr"`
			} `xml:"custom-info"`
		} `xml:"hidden"`
	} `xml:"text_description"`
	Files struct {
		Text string `xml:",chardata"`
		File []struct {
			Text string `xml:",chardata"`
			Size string `xml:"size,attr"`
			Type string `xml:"type,attr"`
		} `xml:"file"`
	} `xml:"files"`
	Categories struct {
		Text      string `xml:",chardata"`
		Categorie []struct {
			Text         string `xml:",chardata"`
			ID           string `xml:"id,attr"`
			CategoryName string `xml:"category_name,attr"`
		} `xml:"categorie"`
	} `xml:"categories"`
	Sequences struct {
		Text     string `xml:",chardata"`
		Sequence []struct {
			Text   string `xml:",chardata"`
			ID     string `xml:"id,attr"`
			Name   string `xml:"name,attr"`
			Number string `xml:"number,attr"`
		} `xml:"sequence"`
	} `xml:"sequences"`
	ArtTags struct {
		Text string `xml:",chardata"`
		Tag  []struct {
			Text     string `xml:",chardata"`
			ID       string `xml:"id,attr"`
			TagTitle string `xml:"tag_title,attr"`
		} `xml:"tag"`
	} `xml:"art_tags"`
	Persons struct {
		Text    string `xml:",chardata"`
		Persons []struct {
			Text   string `xml:",chardata"`
			Role   string `xml:"role,attr"`
			ID     string `xml:"id,attr"`
			Lvl    string `xml:"lvl,attr"`
			First  string `xml:"first,attr"`
			Last   string `xml:"last,attr"`
			Rodit  string `xml:"rodit,attr"`
			Middle string `xml:"middle,attr"`
		} `xml:"persons"`
	} `xml:"persons"`
}
