package sys

import (
	"innodb_inspector/cursor"
	"innodb_inspector/innodb"
)

type FsegEntry struct {
	FsegHdrSpace  uint32 //4
	FsegHdrPageNo uint32 //4
	FsegHdrOffset uint16 //2
}

func NewFsegEntry(data []byte) *FsegEntry {
	br := cursor.New(data)
	return &FsegEntry{
		FsegHdrSpace:  br.Uint32(),
		FsegHdrPageNo: br.Uint32(),
		FsegHdrOffset: br.Uint16(),
	}
}

type DictionaryHeaderPage struct {
	FILHeader *innodb.FILHeader

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

	FILTrailer *innodb.FILTrailer
}

func NewDictionaryHeaderPage(data []byte) *DictionaryHeaderPage {
	br := cursor.New(data)
	br.Seek(38)

	return &DictionaryHeaderPage{
		FILHeader:         innodb.NewFILHeader(data),
		DictHdrRowId:      br.Uint64(),
		DictHdrTableId:    br.Uint64(),
		DictHdrIndexId:    br.Uint64(),
		DictHdrMaxSpaceId: br.Uint32(),
		DictHdrMixIdLow:   br.Uint32(),

		DictHdrTables:   br.Uint32(),
		DictHdrTableIds: br.Uint32(),
		DictHdrColumns:  br.Uint32(),
		DictHdrIndexes:  br.Uint32(),
		DictHdrFields:   br.Uint32(),

		FsegEntry: NewFsegEntry(br.Skip(4).Bytes(10)),

		FILTrailer: innodb.NewFILTrailer(data),
	}
}
