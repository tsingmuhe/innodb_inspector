package xdes

import (
	"innodb_inspector/cursor"
	"innodb_inspector/innodb"
)

const (
	XDESStateFree     XDESState = 1
	XDESStateFreeFrag XDESState = 2
	XDESStateFullFrag XDESState = 3
	XDESStateFseg     XDESState = 4
)

type XDESState uint32

func (t XDESState) String() string {
	switch t {
	case XDESStateFree:
		return "free"
	case XDESStateFreeFrag:
		return "free_frag"
	case XDESStateFullFrag:
		return "full_frag"
	case XDESStateFseg:
		return "fseg"
	}

	return "unknown"
}

func (t XDESState) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

type FSPHeader struct {
	SpaceId       uint32               //4 该文件对应的space id
	Unused        uint32               //4 如其名，保留字节，当前未使用
	Size          uint32               //4 当前表空间总的PAGE个数，扩展文件时需要更新该值
	FreeLimit     uint32               //4 当前尚未初始化的最小Page No。从该Page往后的都尚未加入到表空间的FREE LIST上
	Flags         innodb.Binary        //4 当前表空间的FLAG信息
	FreeFragNUsed uint32               //4 FSP_FREE_FRAG链表上已被使用的Page数，用于快速计算该链表上可用空闲Page数
	Free          *innodb.FlstBaseNode //16 当一个Extent中所有page都未被使用时，放到该链表上，可以用于随后的分配
	FreeFrag      *innodb.FlstBaseNode //16 通常这样的Extent中的Page可能归属于不同的segment
	FullFrag      *innodb.FlstBaseNode //16 Extent中所有的page都被使用掉时，会放到该链表上，当有Page从该Extent释放时，则移回FREE_FRAG链表
	NextSegId     uint64               //8 当前文件中最大Segment ID + 1，用于段分配时的seg id计数器
	FullInodes    *innodb.FlstBaseNode //16 已被完全用满的Inode Page链表
	FreeInodes    *innodb.FlstBaseNode //16 至少存在一个空闲Inode Entry的Inode Page被放到该链表上
}

func NewFSPHeader(data []byte) *FSPHeader {
	br := cursor.New(data)
	br.Seek(38)

	return &FSPHeader{
		SpaceId:       br.Uint32(),
		Unused:        br.Uint32(),
		Size:          br.Uint32(),
		FreeLimit:     br.Uint32(),
		Flags:         br.Bytes(4),
		FreeFragNUsed: br.Uint32(),
		Free:          innodb.NewFlstBaseNode(br.Bytes(16)),
		FreeFrag:      innodb.NewFlstBaseNode(br.Bytes(16)),
		FullFrag:      innodb.NewFlstBaseNode(br.Bytes(16)),
		NextSegId:     br.Uint64(),
		FullInodes:    innodb.NewFlstBaseNode(br.Bytes(16)),
		FreeInodes:    innodb.NewFlstBaseNode(br.Bytes(16)),
	}
}

type XDESEntry struct {
	SegmentId uint64           //8 如果该Extent归属某个segment的话，则记录其ID
	FlstNode  *innodb.FlstNode //12 维持Extent链表的双向指针节点
	State     XDESState        //4 该Extent的状态信息，包括：XDES_FREE，XDES_FREE_FRAG，XDES_FULL_FRAG，XDES_FSEG
	Bitmap    innodb.Binary    //16 总共16*8= 128个bit，用2个bit表示Extent中的一个page，一个bit表示该page是否是空闲的(XDES_FREE_BIT)，另一个保留位，尚未使用（XDES_CLEAN_BIT）
}

func NewXDESEntry(data []byte) []*XDESEntry {
	br := cursor.New(data)
	br.Seek(150)

	var result []*XDESEntry

	for i := 0; i < 256; i++ {
		xdesEntry := &XDESEntry{
			SegmentId: br.Uint64(),
			FlstNode:  innodb.NewFlstNode(br.Bytes(12)),
			State:     XDESState(br.Uint32()),
			Bitmap:    br.Bytes(16),
		}

		if xdesEntry.State > 0 {
			result = append(result, xdesEntry)
		}
	}

	return result
}

type XDESPage struct {
	FILHeader  *innodb.FILHeader
	FSPHeader  *FSPHeader
	XDESEntry  []*XDESEntry
	FILTrailer *innodb.FILTrailer
}

func NewFSPHDRPage(data []byte) *XDESPage {
	return &XDESPage{
		FILHeader:  innodb.NewFILHeader(data),
		FSPHeader:  NewFSPHeader(data),
		XDESEntry:  NewXDESEntry(data),
		FILTrailer: innodb.NewFILTrailer(data),
	}
}

func NewXDESPage(data []byte) *XDESPage {
	return &XDESPage{
		FILHeader:  innodb.NewFILHeader(data),
		XDESEntry:  NewXDESEntry(data),
		FILTrailer: innodb.NewFILTrailer(data),
	}
}
