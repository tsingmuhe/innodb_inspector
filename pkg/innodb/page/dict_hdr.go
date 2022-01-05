package page

const (
	DictionaryHeaderPosition = FilHeaderSize
	DictionaryHeaderSize     = 8 + 8 + 8 + 4 + 4 + 4 + 4 + 4 + 4 + 4 + 4 + 4 + 4 + 2
)

type DictionaryHeader struct {
	buf *Buf

	DictHdrRowId      uint64 //8
	DictHdrTableId    uint64 //8
	DictHdrIndexId    uint64 //8
	DictHdrMaxSpaceId uint32 //4
	DictHdrMixIdLow   uint32 //4

	DictHdrTables   uint32 //4
	DictHdrTableIds uint32 //4
	DictHdrColumns  uint32 //4
	DictHdrIndexes  uint32 //4
	DictHdrFields   uint32 //4

	Unused uint32 `json:"-"`
	
	FsegEntry *FsegEntry
}

func NewDictionaryHeader(pageBytes []byte) *DictionaryHeader {
	from := DictionaryHeaderPosition
	to := from + DictionaryHeaderSize
	buf := NewBuf(from, pageBytes[from:to], to-1)
	return &DictionaryHeader{
		buf:               buf,
		DictHdrRowId:      buf.Uint64(),
		DictHdrTableId:    buf.Uint64(),
		DictHdrIndexId:    buf.Uint64(),
		DictHdrMaxSpaceId: buf.Uint32(),
		DictHdrMixIdLow:   buf.Uint32(),
		DictHdrTables:     buf.Uint32(),
		DictHdrTableIds:   buf.Uint32(),
		DictHdrColumns:    buf.Uint32(),
		DictHdrIndexes:    buf.Uint32(),
		DictHdrFields:     buf.Uint32(),
		Unused:            buf.Uint32(),
		FsegEntry: &FsegEntry{
			FsegHdrSpace:  buf.Uint32(),
			FsegHdrPageNo: buf.Uint32(),
			FsegHdrOffset: buf.Uint16(),
		},
	}
}

func (t *DictionaryHeader) HexEditorTag() *HexEditorTag {
	return &HexEditorTag{
		From:    t.buf.from,
		To:      t.buf.to,
		Color:   "orange",
		Caption: "DictionaryHeader",
	}
}
