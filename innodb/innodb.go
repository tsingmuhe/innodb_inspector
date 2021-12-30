package innodb

import (
	"fmt"
	"innodb_inspector/cursor"
	"strings"
)

const (
	DefaultInnodbPageSize uint = 1024 * 16 //Default InnoDB Page 16K
	UndefinedPageNo            = 4294967295
)

const (
	PageTypeAllocated              PageType = 0
	PageTypeUndoLog                PageType = 2
	PageTypeInode                  PageType = 3
	PageTypeIbufFreeList           PageType = 4
	PageTypeIbufBitmap             PageType = 5
	PageTypeSys                    PageType = 6
	PageTypeTrxSys                 PageType = 7
	PageTypeFspHdr                 PageType = 8
	PageTypeXdes                   PageType = 9
	PageTypeBlob                   PageType = 10
	PageTypeZBlob                  PageType = 11
	PageTypeZBlob2                 PageType = 12
	PageTypeUnknown                PageType = 13
	PageTypeCompressed             PageType = 14
	PageTypeEncrypted              PageType = 15
	PageTypeCompressedAndEncrypted PageType = 16
	PageTypeEncryptedRtree         PageType = 17

	PageTypeIndex PageType = 17855
	PageTypeRtree PageType = 17854
)

type PageType uint16

func (t PageType) String() string {
	switch t {
	case PageTypeAllocated: //Freshly allocated page
		return "FIL_PAGE_TYPE_ALLOCATED"
	case PageTypeUndoLog: //Undo log page
		return "FIL_PAGE_UNDO_LOG"
	case PageTypeInode: //Index node
		return "FIL_PAGE_INODE"
	case PageTypeIbufFreeList: //Insert buffer free list
		return "FIL_PAGE_IBUF_FREE_LIST"
	case PageTypeIbufBitmap: //Insert buffer bitmap
		return "FIL_PAGE_IBUF_BITMAP"
	case PageTypeSys: //System page
		return "FIL_PAGE_TYPE_SYS"
	case PageTypeTrxSys: //Transaction system data
		return "FIL_PAGE_TYPE_TRX_SYS"
	case PageTypeFspHdr: //File space header
		return "FIL_PAGE_TYPE_FSP_HDR"
	case PageTypeXdes: //Extent descriptor page
		return "FIL_PAGE_TYPE_XDES"
	case PageTypeBlob: //Uncompressed BLOB page
		return "FIL_PAGE_TYPE_BLOB"
	case PageTypeZBlob: //First compressed BLOB page
		return "FIL_PAGE_TYPE_ZBLOB"
	case PageTypeZBlob2: //Subsequent compressed BLOB page
		return "FIL_PAGE_TYPE_ZBLOB2"
	case PageTypeUnknown: //In old tablespaces, garbage in FIL_PAGE_TYPE is replaced with this value when flushing pages.
		return "FIL_PAGE_TYPE_UNKNOWN"
	case PageTypeCompressed: //Compressed page
		return "FIL_PAGE_COMPRESSED"
	case PageTypeEncrypted: //Encrypted page
		return "FIL_PAGE_ENCRYPTED"
	case PageTypeCompressedAndEncrypted: //Compressed and Encrypted page
		return "FIL_PAGE_COMPRESSED_AND_ENCRYPTED"
	case PageTypeEncryptedRtree: //Encrypted R-tree page
		return "FIL_PAGE_ENCRYPTED_RTREE"
	case PageTypeIndex: //B-tree node
		return "FIL_PAGE_INDEX"
	case PageTypeRtree: //B-tree node
		return "FIL_PAGE_RTREE"
	}

	return "UNKNOWN"
}

func (t PageType) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func ResolvePageType(data []byte) PageType {
	br := cursor.New(data)
	return PageType(br.Seek(24).Uint16())
}

func ResolveSpaceId(data []byte) uint32 {
	br := cursor.New(data)
	return br.Seek(34).Uint32()
}

//FILHeader 38
type FILHeader struct {
	FilPageSpaceOrChksum      uint32   `json:"FIL_PAGE_SPACE_OR_CHKSUM"`         //4
	FilPageOffset             uint32   `json:"FIL_PAGE_OFFSET"`                  //4相对文件的偏移量
	FilPagePrev               uint32   `json:"FIL_PAGE_PREV"`                    //4
	FilPageNext               uint32   `json:"FIL_PAGE_NEXT"`                    //4
	FilPageLSN                uint64   `json:"FIL_PAGE_LSN"`                     //8
	FilPageType               PageType `json:"FIL_PAGE_TYPE"`                    //2 类型
	FilPageFileFlushLSN       uint64   `json:"FIL_PAGE_FILE_FLUSH_LSN"`          //8
	FilPageArchLogNoOrSpaceId uint32   `json:"FIL_PAGE_ARCH_LOG_NO_OR_SPACE_ID"` //4 表空间id
}

//FILTrailer 8
type FILTrailer struct {
	OldStyleChecksum uint32 `json:"old_style_checksum"` //4
	Low32BitsOfLSN   uint32 `json:"low_32_bits_of_lsn"` //4
}

type PageOverview struct {
	FILHeader  *FILHeader
	FILTrailer *FILTrailer
}

func NewFILHeader(data []byte) *FILHeader {
	br := cursor.New(data)
	return &FILHeader{
		FilPageSpaceOrChksum:      br.Uint32(),
		FilPageOffset:             br.Uint32(),
		FilPagePrev:               br.Uint32(),
		FilPageNext:               br.Uint32(),
		FilPageLSN:                br.Uint64(),
		FilPageType:               PageType(br.Uint16()),
		FilPageFileFlushLSN:       br.Uint64(),
		FilPageArchLogNoOrSpaceId: br.Uint32(),
	}
}

func NewFILTrailer(data []byte) *FILTrailer {
	br := cursor.New(data)
	br.Seek(16376)

	return &FILTrailer{
		OldStyleChecksum: br.Uint32(),
		Low32BitsOfLSN:   br.Uint32(),
	}
}

type Binary []byte

func (t Binary) String() string {
	var elems []string

	for _, i := range t {
		elems = append(elems, fmt.Sprintf("%08b", i))
	}

	return strings.Join(elems, "")
}

func (t Binary) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

//FlstBaseNode 16
type FlstBaseNode struct {
	Len   uint32    //4 存储链表的长度
	First *FlstAddr //6 指向链表的第一个节点
	Last  *FlstAddr //6 指向链表的最后一个节点
}

//FlstNode 12
type FlstNode struct {
	Pre  *FlstAddr //6 指向当前节点的前一个节点
	Next *FlstAddr //6 指向当前节点的下一个节点
}

type FlstAddr struct {
	PageNo uint32 //4 Page No
	Offset uint16 //2 Page内的偏移量
}

func NewFlstBaseNode(data []byte) *FlstBaseNode {
	br := cursor.New(data)
	flstBaseNode := &FlstBaseNode{
		Len: br.Uint32(),
		First: &FlstAddr{
			PageNo: br.Uint32(),
			Offset: br.Uint16(),
		},
		Last: &FlstAddr{
			PageNo: br.Uint32(),
			Offset: br.Uint16(),
		},
	}

	if IsUndefinedPageNo(flstBaseNode.First.PageNo) || IsUndefinedPageNo(flstBaseNode.Last.PageNo) {
		return nil
	}

	return flstBaseNode
}

func NewFlstNode(data []byte) *FlstNode {
	br := cursor.New(data)
	flstNode := &FlstNode{
		Pre: &FlstAddr{
			PageNo: br.Uint32(),
			Offset: br.Uint16(),
		},
		Next: &FlstAddr{
			PageNo: br.Uint32(),
			Offset: br.Uint16(),
		},
	}

	if IsUndefinedPageNo(flstNode.Pre.PageNo) || IsUndefinedPageNo(flstNode.Next.PageNo) {
		return nil
	}

	return flstNode
}

func IsUndefinedPageNo(pageNo uint32) bool {
	return pageNo >= UndefinedPageNo
}
