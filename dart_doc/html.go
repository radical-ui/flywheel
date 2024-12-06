package dart_doc

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type DocHtml struct {
	document *goquery.Document
}

func NewDocHtml(htmlFile string) (*DocHtml, error) {
	file, err := os.Open(htmlFile)
	if err != nil {
		return nil, err
	}

	document, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, err
	}

	return &DocHtml{document}, nil
}

func (self *DocHtml) GetBaseClasses() []string {
	var classes []string

	self.document.Find("ul.clazz-relationships > li > a").Each(func(_ int, selection *goquery.Selection) {
		classes = append(classes, selection.Text())
	})

	return classes
}

func (self *DocHtml) GetFieldTypeData() (string, string, error) {
	selection := self.document.Find("section.multi-line-signature > a").First()

	if selection.Length() == 0 {
		return "", "", fmt.Errorf("failed to find field type signature in html document")
	}

	href, ok := selection.Attr("href")
	if !ok {
		return "", "", fmt.Errorf("no href attribute exists on anchor element")
	}

	return href, selection.Text(), nil
}
