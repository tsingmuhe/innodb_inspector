package page

import (
	"fmt"
	"strings"
)

type Page interface {
	PageNo() uint32
	Type() Type
	Size() int
	SpaceId() uint32
	IsSysTablespace() bool
	FilHeader() *FILHeader
	FILTrailer() *FILTrailer
	Cursor() *Cursor
	String() string
	Notes() string
}

type BasePage struct {
	pageNo   uint32
	pageBits []byte
}

func (f *BasePage) PageNo() uint32 {
	return f.pageNo
}

func (f *BasePage) Type() Type {
	return f.Cursor().Seek(24).Uint16()
}

func (f *BasePage) Size() int {
	return len(f.pageBits) * 8
}

func (f *BasePage) SpaceId() uint32 {
	return f.Cursor().Seek(34).Uint32()
}

func (f *BasePage) IsSysTablespace() bool {
	return f.SpaceId() == 0
}

func (f *BasePage) FilHeader() *FILHeader {
	c := f.Cursor()
	return &FILHeader{
		FilPageSpaceOrChksum:      c.Uint32(),
		FilPageOffset:             c.Uint32(),
		FilPagePrev:               c.Uint32(),
		FilPageNext:               c.Uint32(),
		FilPageLSN:                c.Uint64(),
		FilPageType:               c.Uint16(),
		FilPageFileFlushLSN:       c.Uint64(),
		FilPageArchLogNoOrSpaceId: c.Uint32(),
	}
}

func (f *BasePage) FILTrailer() *FILTrailer {
	c := f.Cursor().Seek(FilTrailerPosition)
	return &FILTrailer{
		OldStyleChecksum: c.Uint32(),
		Low32BitsOfLSN:   c.Uint32(),
	}
}

func (f *BasePage) Cursor() *Cursor {
	return &Cursor{
		data:     f.pageBits,
		position: 0,
	}
}

func (f *BasePage) String() string {
	return ""
}

func (f *BasePage) Notes() string {
	return ""
}

//FlstBaseNode 16
type FlstBaseNode struct {
	Len   uint32   //4 存储链表的长度
	First *Address //6 指向链表的第一个节点
	Last  *Address //6 指向链表的最后一个节点
}

//FlstNode 12
type FlstNode struct {
	Pre  *Address //6 指向当前节点的前一个节点
	Next *Address //6 指向当前节点的下一个节点
}

type Address struct {
	PageNo uint32 //4 Page No
	Offset uint16 //2 Page内的偏移量
}

type Bits []byte

func (t Bits) String() string {
	var elems []string

	for _, i := range t {
		elems = append(elems, fmt.Sprintf("%08b", i))
	}

	return strings.Join(elems, "")
}

func (t Bits) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
