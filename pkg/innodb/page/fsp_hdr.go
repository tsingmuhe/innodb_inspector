package page

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
