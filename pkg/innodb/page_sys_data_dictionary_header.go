package innodb

import (
	"encoding/json"
	"innodb_inspector/pkg/innodb/page"
)

type DictionaryHeaderPage struct {
	*BasePage
}

func (t *DictionaryHeaderPage) DictionaryHeader() *page.DictionaryHeader {
	return page.NewDictionaryHeader(t.pageBytes)
}

func (t *DictionaryHeaderPage) HexEditorTags() []*page.HexEditorTag {
	var result []*page.HexEditorTag
	result = append(result, t.FilHeader().HexEditorTag()...)
	result = append(result, t.DictionaryHeader().HexEditorTag())
	result = append(result, t.FILTrailer().HexEditorTag())
	return result
}

func (t *DictionaryHeaderPage) String() string {
	type Page struct {
		FILHeader        *page.FILHeader
		DictionaryHeader *page.DictionaryHeader
		FILTrailer       *page.FILTrailer
	}

	b, _ := json.MarshalIndent(&Page{
		FILHeader:        t.FilHeader(),
		DictionaryHeader: t.DictionaryHeader(),
		FILTrailer:       t.FILTrailer(),
	}, "", "  ")
	return string(b)
}
