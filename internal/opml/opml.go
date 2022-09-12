package opml

import (
	"encoding/xml"
	"errors"
	"io"
	"time"

	"github.com/indes/flowerss-bot/internal/model"
)

// OPML opml struct
type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

// Head opml head
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

// Body opml body
type Body struct {
	Outlines []Outline `xml:"outline"`
}

// Outline opml outline
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

// NewOPML gen OPML form []byte
func NewOPML(b []byte) (*OPML, error) {
	var root OPML
	err := xml.Unmarshal(b, &root)
	if err != nil {
		return nil, err
	}

	return &root, nil
}

func ReadOPML(r io.Reader) (*OPML, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	o, err := NewOPML(body)
	if err != nil {
		return nil, errors.New("parse opml file error")
	}
	return o, nil
}

// GetFlattenOutlines make all outline at the same xml level
func (o OPML) GetFlattenOutlines() ([]Outline, error) {
	var flattenOutlines []Outline
	for _, line := range o.Body.Outlines {
		if line.Outlines != nil {
			for _, subLine := range line.Outlines {
				// 查找子outline
				if subLine.XMLURL != "" {
					flattenOutlines = append(flattenOutlines, subLine)
				}
			}
		}

		if line.XMLURL != "" {
			flattenOutlines = append(flattenOutlines, line)
		}
	}
	return flattenOutlines, nil
}

// XML dump OPML to xml file
func (o OPML) XML() (string, error) {
	b, err := xml.MarshalIndent(o, "", "\t")
	return xml.Header + string(b), err
}

// ToOPML dump sources to opml file
func ToOPML(sources []*model.Source) (string, error) {
	O := OPML{}
	O.XMLName.Local = "opml"
	O.Version = "2.0"
	O.XMLName.Space = ""
	O.Head.Title = "subscriptions in flowerss"
	O.Head.DateCreated = time.Now().Format(time.RFC1123)
	for _, s := range sources {
		outline := Outline{}
		outline.Text = s.Title
		outline.Type = "rss"
		outline.XMLURL = s.Link
		O.Body.Outlines = append(O.Body.Outlines, outline)
	}
	return O.XML()
}
