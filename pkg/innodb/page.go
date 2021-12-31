package innodb

import "innodb_inspector/pkg/innodb/page"

type Page interface {
	FSPHeaderSpaceId() uint32
	SpaceId() uint32
	PageNo() uint32
	Type() page.Type

	String() string

	IsSysTablespace() bool
	IsDoubleWriteBufferBlock() bool

	Cursor() *Cursor
	CursorAt(position uint32) *Cursor
	CursorAtBodyStart() *Cursor

	FilHeader() *page.FILHeader
	FILTrailer() *page.FILTrailer
}

type BasePage struct {
	fspHeaderSpaceId uint32
	pageNo           uint32
	pageBits         []byte
}

func NewBasePage(fspHeaderSpaceId, pageNo uint32, pageBits []byte) *BasePage {
	return &BasePage{
		fspHeaderSpaceId: fspHeaderSpaceId,
		pageNo:           pageNo,
		pageBits:         pageBits,
	}
}

func (f *BasePage) FSPHeaderSpaceId() uint32 {
	return f.fspHeaderSpaceId
}

func (f *BasePage) SpaceId() uint32 {
	return f.CursorAt(34).Uint32()
}

func (f *BasePage) PageNo() uint32 {
	return f.pageNo
}

func (f *BasePage) Type() page.Type {
	return page.Type(f.CursorAt(24).Uint16())
}

func (f *BasePage) String() string {
	return ""
}

func (f *BasePage) IsSysTablespace() bool {
	return f.fspHeaderSpaceId == 0
}

func (f *BasePage) IsDoubleWriteBufferBlock() bool {
	return f.fspHeaderSpaceId == 0 &&
		f.pageNo >= page.DoubleWriteBufferPageNo1 &&
		f.pageNo < page.DoubleWriteBufferPageNo2
}

func (f *BasePage) Cursor() *Cursor {
	return NewCursor(f.pageBits, 0)
}

func (f *BasePage) CursorAt(position uint32) *Cursor {
	return NewCursor(f.pageBits, position)
}

func (f *BasePage) CursorAtBodyStart() *Cursor {
	return NewCursor(f.pageBits, page.FilHeaderSize)
}

func (f *BasePage) FilHeader() *page.FILHeader {
	c := f.Cursor()
	return &page.FILHeader{
		FilPageSpaceOrChksum:      c.Uint32(),
		FilPageOffset:             c.Uint32(),
		FilPagePrev:               c.Uint32(),
		FilPageNext:               c.Uint32(),
		FilPageLSN:                c.Uint64(),
		FilPageType:               page.Type(c.Uint16()),
		FilPageFileFlushLSN:       c.Uint64(),
		FilPageArchLogNoOrSpaceId: c.Uint32(),
	}
}

func (f *BasePage) FILTrailer() *page.FILTrailer {
	c := f.CursorAt(page.FilTrailerPosition)
	return &page.FILTrailer{
		OldStyleChecksum: c.Uint32(),
		Low32BitsOfLSN:   c.Uint32(),
	}
}
