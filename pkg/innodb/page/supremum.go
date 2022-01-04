package page

type Supremum struct {
	InfoFlags Bits //redundant 7byte compact 5byte
	Supremum  string
}

type CompactSupremum struct {
	OffSet uint32 //112

	InfoFlags        Bits
	RecordType       uint
	NextRecordOffset int16
	Supremum         string
}
