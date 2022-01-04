package page

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

const (
	XDESEntrySize         = 8 + 12 + 4 + 16
	XDESEntryCountPerPage = 256
)

type XDESEntry struct {
	SegmentId uint64    //8 如果该Extent归属某个segment的话，则记录其ID
	FlstNode  *FlstNode //12 维持Extent链表的双向指针节点
	State     XDESState //4 该Extent的状态信息，包括：XDES_FREE，XDES_FREE_FRAG，XDES_FULL_FRAG，XDES_FSEG
	Bitmap    Bits      //16 总共16*8= 128个bit，用2个bit表示Extent中的一个page，一个bit表示该page是否是空闲的(XDES_FREE_BIT)，另一个保留位，尚未使用（XDES_CLEAN_BIT）
}
