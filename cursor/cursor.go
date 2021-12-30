package cursor

import (
	"encoding/binary"
)

type BufferCursor struct {
	data     []byte
	position uint32
}

func (t *BufferCursor) Seek(position uint32) *BufferCursor {
	if position >= 0 {
		t.position = position
	} else {
		t.position = 0
	}
	return t
}

func (t *BufferCursor) Skip(delta uint32) *BufferCursor {
	t.position = t.position + delta
	return t
}

func (t *BufferCursor) Uint16() uint16 {
	return binary.BigEndian.Uint16(t.Bytes(2))
}

func (t *BufferCursor) Uint32() uint32 {
	return binary.BigEndian.Uint32(t.Bytes(4))
}

func (t *BufferCursor) Uint64() uint64 {
	return binary.BigEndian.Uint64(t.Bytes(8))
}

func (t *BufferCursor) Bytes(delta uint32) []byte {
	result := t.data[t.position : t.position+delta]
	t.position = t.position + delta
	return result
}

func New(bytes []byte) *BufferCursor {
	return &BufferCursor{
		data:     bytes,
		position: 0,
	}
}
