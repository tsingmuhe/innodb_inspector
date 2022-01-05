package innodb

import (
	"innodb_inspector/pkg/innodb/page"
	"innodb_inspector/pkg/innodb/tablespace"
)

type Page interface {
	FSPHeaderSpaceId() uint32
	SpaceId() uint32
	PageNo() uint32
	Type() page.Type

	String() string

	HexEditorTags() []*page.HexEditorTag

	IsSysTablespace() bool
	IsDoubleWriteBufferBlock() bool

	PageCursor() *PageCursor
	PageCursorAt(position uint32) *PageCursor
	PageCursorAtBodyStart() *PageCursor

	FilHeader() *page.FILHeader
	FILTrailer() *page.FILTrailer
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

func (f *BasePage) FSPHeaderSpaceId() uint32 {
	return f.fspHeaderSpaceId
}

func (f *BasePage) SpaceId() uint32 {
	return f.PageCursorAt(34).Uint32()
}

func (f *BasePage) PageNo() uint32 {
	return f.pageNo
}

func (f *BasePage) Type() page.Type {
	return page.Type(f.PageCursorAt(24).Uint16())
}

func (f *BasePage) String() string {
	return ""
}

func (f *BasePage) HexEditorTags() []*page.HexEditorTag {
	return nil
}

func (f *BasePage) IsSysTablespace() bool {
	return f.fspHeaderSpaceId == 0
}

func (f *BasePage) IsDoubleWriteBufferBlock() bool {
	return f.fspHeaderSpaceId == 0 &&
		f.pageNo >= tablespace.DoubleWriteBufferPageNo1 &&
		f.pageNo < tablespace.DoubleWriteBufferPageNo2
}

func (f *BasePage) PageCursor() *PageCursor {
	return NewPageCursor(f.pageBytes)
}

func (f *BasePage) PageCursorAt(position uint32) *PageCursor {
	return NewPageCursor(f.pageBytes).SetReaderIndex(position)
}

func (f *BasePage) PageCursorAtBodyStart() *PageCursor {
	return NewPageCursor(f.pageBytes).SetReaderIndex(page.FilHeaderSize)
}

func (f *BasePage) FilHeader() *page.FILHeader {
	return page.NewFILHeader(f.pageBytes)
}

func (f *BasePage) FILTrailer() *page.FILTrailer {
	return page.NewFILTrailer(f.pageBytes)
}

func IsUndefinedPageNo(pageNo uint32) bool {
	return pageNo >= page.UndefinedPageNo
}
