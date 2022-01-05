package innodb

import (
	"encoding/binary"
	"innodb_inspector/pkg/innodb/page"
	"innodb_inspector/pkg/innodb/tablespace"
)

type Page interface {
	FSPHeaderSpaceId() uint32
	SpaceId() uint32
	PageNo() uint32
	Type() page.Type

	IsSysTablespace() bool
	IsDoubleWriteBufferBlock() bool

	FilHeader() *page.FILHeader
	FILTrailer() *page.FILTrailer

	String() string
	HexEditorTags() []*page.HexEditorTag
}

type BasePage struct {
	fspHeaderSpaceId uint32
	pageNo           uint32
	pageBytes        []byte
}

func NewBasePage(fspHeaderSpaceId, pageNo uint32, pageBytes []byte) *BasePage {
	return &BasePage{
		fspHeaderSpaceId: fspHeaderSpaceId,
		pageNo:           pageNo,
		pageBytes:        pageBytes,
	}
}

func (t *BasePage) FSPHeaderSpaceId() uint32 {
	return t.fspHeaderSpaceId
}

func (t *BasePage) SpaceId() uint32 {
	return binary.BigEndian.Uint32(t.pageBytes[34:38])
}

func (t *BasePage) PageNo() uint32 {
	return t.pageNo
}

func (t *BasePage) Type() page.Type {
	b := t.pageBytes[24:26]
	return page.Type(binary.BigEndian.Uint16(b))
}

func (t *BasePage) IsSysTablespace() bool {
	return t.fspHeaderSpaceId == 0
}

func (t *BasePage) IsDoubleWriteBufferBlock() bool {
	return t.fspHeaderSpaceId == 0 &&
		t.pageNo >= tablespace.DoubleWriteBufferPageNo1 &&
		t.pageNo < tablespace.DoubleWriteBufferPageNo2
}

func (t *BasePage) FilHeader() *page.FILHeader {
	return page.NewFILHeader(t.pageBytes)
}

func (t *BasePage) FILTrailer() *page.FILTrailer {
	return page.NewFILTrailer(t.pageBytes)
}

func (t *BasePage) String() string {
	return ""
}

func (t *BasePage) HexEditorTags() []*page.HexEditorTag {
	return make([]*page.HexEditorTag, 0)
}

func IsUndefinedPageNo(pageNo uint32) bool {
	return pageNo >= page.UndefinedPageNo
}
