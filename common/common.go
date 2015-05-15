package common

import (
	"github.com/biogo/store/interval"
	"fmt"
	"io"
	"github.com/golang/protobuf/proto"
	"encoding/binary"
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

func WriteMessage(w io.Writer, m proto.Message) error {
	binaryMessage, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.LittleEndian, int32(len(binaryMessage)))
	if err != nil {
		return err
	}
	_, err = w.Write(binaryMessage)
	return err
}

func ReadMessage(r io.Reader, m proto.Message) (err error) {
	var messageLen int32
	err = binary.Read(r, binary.LittleEndian, &messageLen)
	if err != nil {
		return err
	}
	messageBytes := make([]byte, messageLen)
	_, err = r.Read(messageBytes)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(messageBytes, m)
	return err
}