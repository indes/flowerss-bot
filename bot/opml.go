package bot

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

type Head struct {
	Title           string `xml:"title"`
	DateCreated     string `xml:"dateCreated,omitempty"`
	DateModified    string `xml:"dateModified,omitempty"`
	OwnerName       string `xml:"ownerName,omitempty"`
	OwnerEmail      string `xml:"ownerEmail,omitempty"`
	OwnerID         string `xml:"ownerId,omitempty"`
	Docs            string `xml:"docs,omitempty"`
	ExpansionState  string `xml:"expansionState,omitempty"`
	VertScrollState string `xml:"vertScrollState,omitempty"`
	WindowTop       string `xml:"windowTop,omitempty"`
	WindowBottom    string `xml:"windowBottom,omitempty"`
	WindowLeft      string `xml:"windowLeft,omitempty"`
	WindowRight     string `xml:"windowRight,omitempty"`
}

type Body struct {
	Outlines []Outline `xml:"outline"`
}

type Outline struct {
	Outlines     []Outline `xml:"outline"`
	Text         string    `xml:"text,attr"`
	Type         string    `xml:"type,attr,omitempty"`
	IsComment    string    `xml:"isComment,attr,omitempty"`
	IsBreakpoint string    `xml:"isBreakpoint,attr,omitempty"`
	Created      string    `xml:"created,attr,omitempty"`
	Category     string    `xml:"category,attr,omitempty"`
	XMLURL       string    `xml:"xmlUrl,attr,omitempty"`
	HTMLURL      string    `xml:"htmlUrl,attr,omitempty"`
	URL          string    `xml:"url,attr,omitempty"`
	Language     string    `xml:"language,attr,omitempty"`
	Title        string    `xml:"title,attr,omitempty"`
	Version      string    `xml:"version,attr,omitempty"`
	Description  string    `xml:"description,attr,omitempty"`
}

func NewOPML(b []byte) (*OPML, error) {
	var root OPML
	err := xml.Unmarshal(b, &root)
	if err != nil {
		return nil, err
	}

	return &root, nil
}

func GetOPMLByURL(file_url string) (*OPML, error) {

	proxy, _ := url.Parse("socks5://" + socks5Proxy)
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5,
	}
	resp, err := client.Get(file_url)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	o, err := NewOPML(body)
	fmt.Println(string(body))

	if err != nil {
		return nil, err
	}
	return o, err
}

func (doc OPML) Outlines() []Outline {
	return doc.Body.Outlines
}

func (doc OPML) XML() (string, error) {
	b, err := xml.MarshalIndent(doc, "", "\t")
	return xml.Header + string(b), err
}
