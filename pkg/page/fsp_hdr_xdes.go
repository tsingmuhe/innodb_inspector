package page

import (
	"encoding/json"
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
	SpaceId       uint32        //4 该文件对应的space id
	Unused        uint32        //4 如其名，保留字节，当前未使用
	Size          uint32        //4 当前表空间总的PAGE个数，扩展文件时需要更新该值
	FreeLimit     uint32        //4 当前尚未初始化的最小Page No。从该Page往后的都尚未加入到表空间的FREE LIST上
	Flags         Bits          //4 当前表空间的FLAG信息
	FreeFragNUsed uint32        //4 FSP_FREE_FRAG链表上已被使用的Page数，用于快速计算该链表上可用空闲Page数
	Free          *FlstBaseNode //16 当一个Extent中所有page都未被使用时，放到该链表上，可以用于随后的分配
	FreeFrag      *FlstBaseNode //16 通常这样的Extent中的Page可能归属于不同的segment
	FullFrag      *FlstBaseNode //16 Extent中所有的page都被使用掉时，会放到该链表上，当有Page从该Extent释放时，则移回FREE_FRAG链表
	NextSegId     uint64        //8 当前文件中最大Segment ID + 1，用于段分配时的seg id计数器
	FullInodes    *FlstBaseNode //16 已被完全用满的Inode Page链表
	FreeInodes    *FlstBaseNode //16 至少存在一个空闲Inode Entry的Inode Page被放到该链表上
}

type XDESEntry struct {
	SegmentId uint64    //8 如果该Extent归属某个segment的话，则记录其ID
	FlstNode  *FlstNode //12 维持Extent链表的双向指针节点
	State     XDESState //4 该Extent的状态信息，包括：XDES_FREE，XDES_FREE_FRAG，XDES_FULL_FRAG，XDES_FSEG
	Bitmap    Bits      //16 总共16*8= 128个bit，用2个bit表示Extent中的一个page，一个bit表示该page是否是空闲的(XDES_FREE_BIT)，另一个保留位，尚未使用（XDES_CLEAN_BIT）
}

type FspHdrXdesPage struct {
	*BasePage
}

func (f *FspHdrXdesPage) IsFspHdr() bool {
	return f.pageNo == 0
}

func (f *FspHdrXdesPage) FSPHeader() *FSPHeader {
	if !f.IsFspHdr() {
		return nil
	}

	c := f.CursorAtBodyStart()
	return &FSPHeader{
		SpaceId:       c.Uint32(),
		Unused:        c.Uint32(),
		Size:          c.Uint32(),
		FreeLimit:     c.Uint32(),
		Flags:         c.Bytes(4),
		FreeFragNUsed: c.Uint32(),
		Free:          c.FlstBaseNode(),
		FreeFrag:      c.FlstBaseNode(),
		FullFrag:      c.FlstBaseNode(),
		NextSegId:     c.Uint64(),
		FullInodes:    c.FlstBaseNode(),
		FreeInodes:    c.FlstBaseNode(),
	}
}

func (f *FspHdrXdesPage) XDESEntry() []*XDESEntry {
	c := f.CursorAt(150)

	var xdesEntries []*XDESEntry

	for i := 0; i < 256; i++ {
		xdesEntry := &XDESEntry{
			SegmentId: c.Uint64(),
			FlstNode:  c.FlstNode(),
			State:     XDESState(c.Uint32()),
			Bitmap:    c.Bytes(16),
		}

		if xdesEntry.State > 0 {
			xdesEntries = append(xdesEntries, xdesEntry)
		}
	}

	return xdesEntries
}

func (f *FspHdrXdesPage) String() string {
	type Page struct {
		FILHeader   *FILHeader
		FSPHeader   *FSPHeader
		XDESEntries []*XDESEntry
		FILTrailer  *FILTrailer
	}

	b, _ := json.MarshalIndent(&Page{
		FILHeader:   f.FilHeader(),
		FSPHeader:   f.FSPHeader(),
		XDESEntries: f.XDESEntry(),
		FILTrailer:  f.FILTrailer(),
	}, "", "  ")
	return string(b)
}
