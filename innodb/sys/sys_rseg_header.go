package sys

import (
	"innodb_inspector/cursor"
	"innodb_inspector/innodb"
)

type SysRsegHeaderPage struct {
	FILHeader   *innodb.FILHeader
	MaxSize     uint32
	HistorySize uint32
	HistoryList *innodb.FlstBaseNode
	FsegEntry   *FsegEntry
	FILTrailer  *innodb.FILTrailer
}

func NewSysRsegHeaderPage(data []byte) *SysRsegHeaderPage {
	br := cursor.New(data)
	br.Seek(38)
	return &SysRsegHeaderPage{
		FILHeader:   innodb.NewFILHeader(data),
		MaxSize:     br.Uint32(),
		HistorySize: br.Uint32(),
		HistoryList: innodb.NewFlstBaseNode(br.Bytes(16)),
		FsegEntry:   NewFsegEntry(br.Bytes(10)),
		FILTrailer:  innodb.NewFILTrailer(data),
	}
}
