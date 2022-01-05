package page

const (
	UndefinedPageNo = 4294967295

	DefaultSize uint32 = 1024 * 16 //Default InnoDB Page 16K
)

const (
	InsertBufferHeaderPageNo      = 3
	InsertBufferRootPageNo        = 4
	TransactionSystemHeaderPageNo = 5
	FirstRollbackSegmentPageNo    = 6
	DataDictionaryHeaderPageNo    = 7

	DoubleWriteBufferPageNo1 = 64
	DoubleWriteBufferPageNo2 = 192

	RootPageOfFirstIndexPageNo  = 3
	RootPageOfSecondIndexPageNo = 4
)

const (
	FilHeaderPosition = 0
	FilHeaderSize     = 4 + 4 + 4 + 4 + 8 + 2 + 8 + 4
)

type FILHeader struct {
	FilPageSpaceOrChksum      uint32
	FilPageOffset             uint32
	FilPagePrev               uint32
	FilPageNext               uint32
	FilPageLSN                uint64
	FilPageType               Type
	FilPageFileFlushLSN       uint64
	FilPageArchLogNoOrSpaceId uint32
}

func (t *FILHeader) HexEditorTag() *HexEditorTag {
	return &HexEditorTag{
		From:    0,
		To:      FilHeaderSize - 1,
		Color:   "red",
		Caption: "FILHeader",
	}
}

const (
	FilTrailerSize uint32 = 4 + 4
)

type FILTrailer struct {
	OldStyleChecksum uint32
	Low32BitsOfLSN   uint32
}

func (t *FILTrailer) HexEditorTag(pageSize int) *HexEditorTag {
	return &HexEditorTag{
		From:    uint32(pageSize) - FilTrailerSize,
		To:      uint32(pageSize) - 1,
		Color:   "red",
		Caption: "FILTrailer",
	}
}
