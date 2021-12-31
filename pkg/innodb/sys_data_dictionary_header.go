package innodb

import (
	"encoding/json"
	"innodb_inspector/pkg/innodb/page"
)

type DictionaryHeaderPage struct {
	*BasePage
}

func (t *DictionaryHeaderPage) DictionaryHeader() *page.DictionaryHeader {
	c := t.CursorAtBodyStart()
	dictionaryHeader := &page.DictionaryHeader{
		DictHdrRowId:      c.Uint64(),
		DictHdrTableId:    c.Uint64(),
		DictHdrIndexId:    c.Uint64(),
		DictHdrMaxSpaceId: c.Uint32(),
		DictHdrMixIdLow:   c.Uint32(),
		DictHdrTables:     c.Uint32(),
		DictHdrTableIds:   c.Uint32(),
		DictHdrColumns:    c.Uint32(),
		DictHdrIndexes:    c.Uint32(),
		DictHdrFields:     c.Uint32(),
	}

	c.Skip(4)

	dictionaryHeader.FsegEntry = &page.FsegEntry{
		FsegHdrSpace:  c.Uint32(),
		FsegHdrPageNo: c.Uint32(),
		FsegHdrOffset: c.Uint16(),
	}

	return dictionaryHeader
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
