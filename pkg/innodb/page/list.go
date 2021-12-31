package page

//FlstBaseNode 16
type FlstBaseNode struct {
	Len   uint32   //4 存储链表的长度
	First *Address //6 指向链表的第一个节点
	Last  *Address //6 指向链表的最后一个节点
}

//FlstNode 12
type FlstNode struct {
	Pre  *Address //6 指向当前节点的前一个节点
	Next *Address //6 指向当前节点的下一个节点
}

type Address struct {
	PageNo uint32 //4 Page No
	Offset uint16 //2 Page内的偏移量
}
