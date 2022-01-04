package page

type Infimum struct {
	InfoFlags Bits //redundant 7byte compact 5byte
	Infimum   string
}

const (
	CompactInfimumPosition = FilHeaderSize + IndexHeaderSize + FSegHeaderSize
	CompactInfimumSize     = 3 + 2 + 8
)

type CompactInfimum struct {
	OffSet uint32 //99

	InfoFlags        Bits
	RecordType       uint
	NextRecordOffset int16
	Infimum          string
}

func (t *CompactInfimum) NextRecord() uint32 {
	return uint32(int32(t.OffSet) + int32(t.NextRecordOffset))
}
