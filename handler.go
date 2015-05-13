package main
import (
    "net"
    "log"
    "bufio"
    "github.com/bo0rsh201/intervals/common"
    "bytes"
    "fmt"
	"github.com/golang/protobuf/proto"
	"github.com/bo0rsh201/intervals/proto"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()

	var request = new(messages.IntervalRequest)

	bufSlice := make([]byte, 4096)
	buf := proto.NewBuffer(bufSlice)

	err := buf.Unmarshal(request)
	if err != nil {
		log.Print(err)
		return
	}

    var debugBuffer bytes.Buffer
    debugBuffer.WriteString(fmt.Sprintf("Got request for %d\nResults:\n", request.Point))

    mutex.RLock()
    matches := Data.Get(common.IntInterval{Start: int(*request.Point), End: int(*request.Point)})
    mutex.RUnlock()

	var response = new(messages.IntervalResponse)
    for _, match := range matches {
		response.Points = append(response.Points, int32(match.ID()))
        debugBuffer.WriteString(fmt.Sprintf("Id: %d Range: %d - %d\n", match.ID(), match.Range().Start, match.Range().End))
    }
    debugBuffer.WriteString("\n")

	err = buf.Marshal(response)
    if err != nil {
        log.Print(err)
        return
    }

	writer := bufio.NewWriter(conn)

	_, err = writer.Write(buf.Bytes())
	if err != nil {
		log.Print(err)
	}

    fmt.Print(debugBuffer.String())
}
