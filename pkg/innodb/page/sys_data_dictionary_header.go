package page

import "encoding/json"

type DictionaryHeader struct {
	DictHdrRowId      uint64 //8
	DictHdrTableId    uint64 //8
	DictHdrIndexId    uint64 //8
	DictHdrMaxSpaceId uint32 //4
	DictHdrMixIdLow   uint32 //4

	DictHdrTables   uint32 //4
	DictHdrTableIds uint32 //4
	DictHdrColumns  uint32 //4
	DictHdrIndexes  uint32 //4
	DictHdrFields   uint32 //4

	//skip 4
	FsegEntry *FsegEntry
}

type DictionaryHeaderPage struct {
	*BasePage
}

func (t *DictionaryHeaderPage) DictionaryHeader() *DictionaryHeader {
	c := t.CursorAtBodyStart()
	dictionaryHeader := &DictionaryHeader{
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

	dictionaryHeader.FsegEntry = &FsegEntry{
		FsegHdrSpace:  c.Uint32(),
		FsegHdrPageNo: c.Uint32(),
		FsegHdrOffset: c.Uint16(),
	}

	return dictionaryHeader
}

func (t *DictionaryHeaderPage) String() string {
	type Page struct {
		FILHeader        *FILHeader
		DictionaryHeader *DictionaryHeader
		FILTrailer       *FILTrailer
	}

	b, _ := json.MarshalIndent(&Page{
		FILHeader:        t.FilHeader(),
		DictionaryHeader: t.DictionaryHeader(),
		FILTrailer:       t.FILTrailer(),
	}, "", "  ")
	return string(b)
}
