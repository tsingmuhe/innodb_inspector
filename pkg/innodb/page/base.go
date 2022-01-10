package page

const (
	UndefinedPageNo        = 4294967295
	DefaultSize     uint32 = 1024 * 16 //Default InnoDB Page 16K
)

const (
	FilHeaderPosition uint32 = 0
	FilHeaderSize     uint32 = 4 + 4 + 4 + 4 + 8 + 2 + 8 + 4
)

type FILHeader struct {
	buf *Buf

	FilPageSpaceOrChksum      uint32
	FilPageOffset             uint32
	FilPagePrev               uint32
	FilPageNext               uint32
	FilPageLSN                uint64
	FilPageType               Type
	FilPageFileFlushLSN       uint64
	FilPageArchLogNoOrSpaceId uint32
}

func NewFILHeader(pageBytes []byte) *FILHeader {
	from := FilHeaderPosition
	to := FilHeaderPosition + FilHeaderSize

	buf := NewBuf(from, pageBytes[from:to], to-1)
	return &FILHeader{
		buf:                       buf,
		FilPageSpaceOrChksum:      buf.Uint32(),
		FilPageOffset:             buf.Uint32(),
		FilPagePrev:               buf.Uint32(),
		FilPageNext:               buf.Uint32(),
		FilPageLSN:                buf.Uint64(),
		FilPageType:               Type(buf.Uint16()),
		FilPageFileFlushLSN:       buf.Uint64(),
		FilPageArchLogNoOrSpaceId: buf.Uint32(),
	}
}

func (t *FILHeader) HexEditorTag() []*HexEditorTag {
	var tags []*HexEditorTag
	tags = append(tags, &HexEditorTag{
		From:    t.buf.from,
		To:      t.buf.from + 23,
		Color:   "red",
		Caption: "FILHeader",
	}, &HexEditorTag{
		From:    t.buf.from + 24,
		To:      t.buf.from + 25,
		Color:   "green",
		Caption: "FILHeader",
	}, &HexEditorTag{
		From:    t.buf.from + 26,
		To:      t.buf.to,
		Color:   "red",
		Caption: "FILHeader",
	})
	return tags
}

const (
	FilTrailerSize uint32 = 4 + 4
)

type FILTrailer struct {
	buf *Buf

	OldStyleChecksum uint32
	Low32BitsOfLSN   uint32
}

func NewFILTrailer(pageBytes []byte) *FILTrailer {
	pageSize := uint32(len(pageBytes))
	from := pageSize - FilTrailerSize
	to := pageSize

	buf := NewBuf(from, pageBytes[from:to], to-1)
	return &FILTrailer{
		buf:              buf,
		OldStyleChecksum: buf.Uint32(),
		Low32BitsOfLSN:   buf.Uint32(),
	}
}

func (t *FILTrailer) HexEditorTag() *HexEditorTag {
	return &HexEditorTag{
		From:    t.buf.from,
		To:      t.buf.to,
		Color:   "red",
		Caption: "FILTrailer",
	}
}

const (
	FlstBaseNodeSize = 4 + FlstAddressSize + FlstAddressSize
)

type FlstBaseNode struct {
	Len   uint32       //4 存储链表的长度
	First *FlstAddress //6 指向链表的第一个节点
	Last  *FlstAddress //6 指向链表的最后一个节点
}

const (
	FlstNodeSize = FlstAddressSize + FlstAddressSize
)

type FlstNode struct {
	Pre  *FlstAddress //6 指向当前节点的前一个节点
	Next *FlstAddress //6 指向当前节点的下一个节点
}

func NewFlstNode(from uint32, pageBytes []byte) *FlstNode {
	to := from + FlstNodeSize
	buf := NewBuf(from, pageBytes[from:to], to-1)
	return &FlstNode{
		Pre: &FlstAddress{
			PageNo: buf.Uint32(),
			Offset: buf.Uint16(),
		},
		Next: &FlstAddress{
			PageNo: buf.Uint32(),
			Offset: buf.Uint16(),
		},
	}
}

const (
	FlstAddressSize = 4 + 2
)

type FlstAddress struct {
	PageNo uint32 //4 Page No
	Offset uint16 //2 Page内的偏移量
}

type HexEditorTag struct {
	From    uint32 `json:"from"`
	To      uint32 `json:"to"`
	Color   string `json:"color"`
	Caption string `json:"caption"`
}
