package page

type Supremum struct {
	InfoFlags Bits //redundant 7byte compact 5byte
	Supremum  string
}

const (
	CompactSupremumPosition = FilHeaderSize + IndexHeaderSize + FSegHeaderSize + CompactInfimumSize
	CompactSupremumSize     = 3 + 2 + 8
)

type CompactSupremum struct {
	OffSet uint32 //112

	InfoFlags        Bits
	RecordType       uint
	NextRecordOffset int16
	Supremum         string
}
