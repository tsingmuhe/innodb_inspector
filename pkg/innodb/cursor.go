package innodb

import (
	"encoding/binary"
	"innodb_inspector/pkg/innodb/page"
)

type PageCursor struct {
	data        []byte
	readerIndex uint32
}

func NewPageCursor(data []byte) *PageCursor {
	return &PageCursor{
		data:        data,
		readerIndex: 0,
	}
}

func (t *PageCursor) ReaderIndex() uint32 {
	return t.readerIndex
}

func (t *PageCursor) SetReaderIndex(readerIndex uint32) *PageCursor {
	t.readerIndex = readerIndex
	return t
}

func (t *PageCursor) SkipBytes(length uint32) *PageCursor {
	t.readerIndex = t.readerIndex + length
	return t
}

func (t *PageCursor) Byte() byte {
	result := t.data[t.readerIndex]
	t.readerIndex = t.readerIndex + 1
	return result
}

func (t *PageCursor) Bytes(length uint32) []byte {
	result := t.data[t.readerIndex : t.readerIndex+length]
	t.readerIndex = t.readerIndex + length
	return result
}

func (t *PageCursor) Uint8() uint8 {
	return uint8(t.Byte())
}

func (t *PageCursor) Int8() int8 {
	return int8(t.Byte())
}

func (t *PageCursor) Uint16() uint16 {
	return binary.BigEndian.Uint16(t.Bytes(2))
}

func (t *PageCursor) Int16() int16 {
	return int16(t.Uint16())
}

func (t *PageCursor) Uint32() uint32 {
	return binary.BigEndian.Uint32(t.Bytes(4))
}

func (t *PageCursor) Int32() int32 {
	return int32(t.Uint32())
}

func (t *PageCursor) Uint64() uint64 {
	return binary.BigEndian.Uint64(t.Bytes(8))
}

func (t *PageCursor) Int64() int64 {
	return int64(t.Uint64())
}

func (t *PageCursor) FlstBaseNode() *page.FlstBaseNode {
	return &page.FlstBaseNode{
		Len: t.Uint32(),
		First: &page.FlstAddress{
			PageNo: t.Uint32(),
			Offset: t.Uint16(),
		},
		Last: &page.FlstAddress{
			PageNo: t.Uint32(),
			Offset: t.Uint16(),
		},
	}
}
