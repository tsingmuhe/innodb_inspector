package page

const (
	IndexHeaderPosition = FilHeaderSize
	IndexHeaderSize     = 2 + 2 + 2 + 2 + 2 + 2 + 2 + 2 + 2 + 8 + 2 + 8
)

type IndexHeader struct {
	NDirSlots  uint16
	HeapTop    uint16
	NHeap      uint16
	Free       uint16
	Garbage    uint16
	LastInsert uint16
	Direction  uint16
	NDirection uint16
	NRecs      uint16
	MaxTrxID   uint64
	Level      uint16
	IndexId    uint64
}
