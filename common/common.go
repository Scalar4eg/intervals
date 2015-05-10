package common

import (
	"github.com/biogo/store/interval"
	"fmt"
)

type IntInterval struct {
	Start, End int
	Id uintptr
}

func (i IntInterval) Overlap(b interval.IntRange) bool {
	return i.End > b.Start && i.Start < b.End
}
func (i IntInterval) ID() uintptr {
	return i.Id
}
func (i IntInterval) Range() interval.IntRange {
	return interval.IntRange{i.Start, i.End}
}
func (i IntInterval) String() string {
	return fmt.Sprintf("[%d,%d)#%d", i.Start, i.End, i.ID)
}