package page

const (
	FsegEntrySize = 4 + 4 + 2
)

type FsegEntry struct {
	FsegHdrSpace  uint32 //4
	FsegHdrPageNo uint32 //4
	FsegHdrOffset uint16 //2
}
