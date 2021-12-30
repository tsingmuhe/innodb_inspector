package page

const (
	FilPageTypeAllocated          Type = 0
	FilPageUndoLog                Type = 2
	FilPageInode                  Type = 3
	FilPageIbufFreeList           Type = 4
	FilPageIbufBitmap             Type = 5
	FilPageTypeSys                Type = 6
	FilPageTypeTrxSys             Type = 7
	FilPageTypeFspHdr             Type = 8
	FilPageTypeXdes               Type = 9
	FilPageTypeBlob               Type = 10
	FilPageTypeZblob              Type = 11
	FilPageTypeZblob2             Type = 12
	FilPageTypeUnknown            Type = 13
	FilPageCompressed             Type = 14
	FilPageEncrypted              Type = 15
	FilPageCompressedAndEncrypted Type = 16
	FilPageEncryptedRtree         Type = 17

	FilPageIndex Type = 17855
	FilPageRtree Type = 17854
)

type Type uint16

func (t Type) String() string {
	switch t {
	case FilPageTypeAllocated:
		return "FIL_PAGE_TYPE_ALLOCATED"
	case FilPageUndoLog:
		return "FIL_PAGE_UNDO_LOG"
	case FilPageInode:
		return "FIL_PAGE_INODE"
	case FilPageIbufFreeList:
		return "FIL_PAGE_IBUF_FREE_LIST"
	case FilPageIbufBitmap:
		return "FIL_PAGE_IBUF_BITMAP"
	case FilPageTypeSys:
		return "FIL_PAGE_TYPE_SYS"
	case FilPageTypeTrxSys:
		return "FIL_PAGE_TYPE_TRX_SYS"
	case FilPageTypeFspHdr:
		return "FIL_PAGE_TYPE_FSP_HDR"
	case FilPageTypeXdes:
		return "FIL_PAGE_TYPE_XDES"
	case FilPageTypeBlob:
		return "FIL_PAGE_TYPE_BLOB"
	case FilPageTypeZblob:
		return "FIL_PAGE_TYPE_ZBLOB"
	case FilPageTypeZblob2:
		return "FIL_PAGE_TYPE_ZBLOB2"
	case FilPageTypeUnknown:
		return "FIL_PAGE_TYPE_UNKNOWN"
	case FilPageCompressed:
		return "FIL_PAGE_COMPRESSED"
	case FilPageEncrypted:
		return "FIL_PAGE_ENCRYPTED"
	case FilPageCompressedAndEncrypted:
		return "FIL_PAGE_COMPRESSED_AND_ENCRYPTED"
	case FilPageEncryptedRtree:
		return "FIL_PAGE_ENCRYPTED_RTREE"
	case FilPageIndex:
		return "FIL_PAGE_INDEX"
	case FilPageRtree:
		return "FIL_PAGE_RTREE"
	}

	return "unknown page type"
}

func (t Type) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
