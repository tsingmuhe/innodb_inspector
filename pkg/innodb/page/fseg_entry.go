package page

type FsegEntry struct {
	FsegHdrSpace  uint32 //4
	FsegHdrPageNo uint32 //4
	FsegHdrOffset uint16 //2
}
