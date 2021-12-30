package sys

import (
	"innodb_inspector/cursor"
	"innodb_inspector/innodb"
)

type IBufHeaderPage struct {
	FILHeader  *innodb.FILHeader
	FsegEntry  *FsegEntry
	FILTrailer *innodb.FILTrailer
}

func NewIBufHeaderPage(data []byte) *IBufHeaderPage {
	return &IBufHeaderPage{
		FILHeader:  innodb.NewFILHeader(data),
		FsegEntry:  NewFsegEntry(cursor.New(data).Seek(38).Bytes(10)),
		FILTrailer: innodb.NewFILTrailer(data),
	}
}
