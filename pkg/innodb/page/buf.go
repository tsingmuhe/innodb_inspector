package page

import "encoding/binary"

type Buf struct {
	from uint32
	to   uint32

	bytes     []byte
	readIndex uint32
}

func NewBuf(from uint32, bytes []byte, to uint32) *Buf {
	return &Buf{
		from:      from,
		to:        to,
		bytes:     bytes,
		readIndex: 0,
	}
}

func (t *Buf) SkipBytes(length uint32) *Buf {
	t.readIndex = t.readIndex + length
	return t
}

func (t *Buf) Byte() byte {
	result := t.bytes[t.readIndex]
	t.readIndex = t.readIndex + 1
	return result
}

func (t *Buf) Bytes(length uint32) []byte {
	result := t.bytes[t.readIndex : t.readIndex+length]
	t.readIndex = t.readIndex + length
	return result
}

func (t *Buf) Uint8() uint8 {
	return uint8(t.Byte())
}

func (t *Buf) Int8() int8 {
	return int8(t.Byte())
}

func (t *Buf) Uint16() uint16 {
	return binary.BigEndian.Uint16(t.Bytes(2))
}

func (t *Buf) Int16() int16 {
	return int16(t.Uint16())
}

func (t *Buf) Uint32() uint32 {
	return binary.BigEndian.Uint32(t.Bytes(4))
}

func (t *Buf) Int32() int32 {
	return int32(t.Uint32())
}

func (t *Buf) Uint64() uint64 {
	return binary.BigEndian.Uint64(t.Bytes(8))
}

func (t *Buf) Int64() int64 {
	return int64(t.Uint64())
}

func (t *Buf) FlstBaseNode() *FlstBaseNode {
	return &FlstBaseNode{
		Len: t.Uint32(),
		First: &FlstAddress{
			PageNo: t.Uint32(),
			Offset: t.Uint16(),
		},
		Last: &FlstAddress{
			PageNo: t.Uint32(),
			Offset: t.Uint16(),
		},
	}
}

func (t *Buf) FlstNode() *FlstNode {
	return &FlstNode{
		Pre: &FlstAddress{
			PageNo: t.Uint32(),
			Offset: t.Uint16(),
		},
		Next: &FlstAddress{
			PageNo: t.Uint32(),
			Offset: t.Uint16(),
		},
	}
}
