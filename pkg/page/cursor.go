package page

import (
	"encoding/binary"
)

type Cursor struct {
	data     []byte
	position uint32
}

func (t *Cursor) Skip(delta uint32) *Cursor {
	t.position = t.position + delta
	return t
}

func (t *Cursor) Uint16() uint16 {
	return binary.BigEndian.Uint16(t.Bytes(2))
}

func (t *Cursor) Uint32() uint32 {
	return binary.BigEndian.Uint32(t.Bytes(4))
}

func (t *Cursor) Uint64() uint64 {
	return binary.BigEndian.Uint64(t.Bytes(8))
}

func (t *Cursor) Bytes(delta uint32) []byte {
	result := t.data[t.position : t.position+delta]
	t.position = t.position + delta
	return result
}

func (t *Cursor) FlstBaseNode() *FlstBaseNode {
	return &FlstBaseNode{
		Len: t.Uint32(),
		First: &Address{
			PageNo: t.Uint32(),
			Offset: t.Uint16(),
		},
		Last: &Address{
			PageNo: t.Uint32(),
			Offset: t.Uint16(),
		},
	}
}

func (t *Cursor) FlstNode() *FlstNode {
	return &FlstNode{
		Pre: &Address{
			PageNo: t.Uint32(),
			Offset: t.Uint16(),
		},
		Next: &Address{
			PageNo: t.Uint32(),
			Offset: t.Uint16(),
		},
	}
}
