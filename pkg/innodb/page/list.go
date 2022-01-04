package page

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

const (
	FlstAddressSize = 4 + 2
)

type FlstAddress struct {
	PageNo uint32 //4 Page No
	Offset uint16 //2 Page内的偏移量
}
