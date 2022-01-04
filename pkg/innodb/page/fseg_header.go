package page

const (
	FSegHeaderPosition = FilHeaderSize + IndexHeaderSize
	FSegHeaderSize     = FsegEntrySize + FsegEntrySize
)

type FSegHeader struct {
	Leaf   *FsegEntry
	NoLeaf *FsegEntry
}
