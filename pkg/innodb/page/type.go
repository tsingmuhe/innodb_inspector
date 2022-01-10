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
		return "Freshly allocated page"
	case FilPageUndoLog:
		return "Undo log page"
	case FilPageInode:
		return "Inode page"
	case FilPageIbufFreeList:
		return "Insert buffer free list page"
	case FilPageIbufBitmap:
		return "Insert buffer bitmap"
	case FilPageTypeSys:
		return "System page"
	case FilPageTypeTrxSys:
		return "Transaction system page"
	case FilPageTypeFspHdr:
		return "File Space Header"
	case FilPageTypeXdes:
		return "Extent descriptor page"
	case FilPageTypeBlob:
		return "BLOB page"
	case FilPageTypeZblob:
		return "FIL_PAGE_TYPE_ZBLOB"
	case FilPageTypeZblob2:
		return "FIL_PAGE_TYPE_ZBLOB2"
	case FilPageTypeUnknown:
		return "FIL_PAGE_TYPE_UNKNOWN"
	case FilPageCompressed:
		return "Compressed BLOB page"
	case FilPageEncrypted:
		return "FIL_PAGE_ENCRYPTED"
	case FilPageCompressedAndEncrypted:
		return "FIL_PAGE_COMPRESSED_AND_ENCRYPTED"
	case FilPageEncryptedRtree:
		return "FIL_PAGE_ENCRYPTED_RTREE"
	case FilPageIndex:
		return "Index page"
	case FilPageRtree:
		return "FIL_PAGE_RTREE"
	}

	return "Other type of page"
}

func (t Type) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
