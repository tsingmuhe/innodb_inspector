package page

const (
	FilHeaderPosition  = 0
	FilTrailerPosition = 0

	FilHeaderSize  = 4 + 4 + 4 + 4 + 8 + 2 + 8 + 4
	FilTrailerSize = 4 + 4
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

type FILTrailer struct {
	OldStyleChecksum uint32
	Low32BitsOfLSN   uint32
}
