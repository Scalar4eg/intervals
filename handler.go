package main
import (
    "net"
    "log"
    "github.com/bo0rsh201/intervals/common"
    "fmt"
	"github.com/bo0rsh201/intervals/proto"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()

	var request = messages.IntervalRequest{}

	err := common.ReadMessage(conn, &request)
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Printf("Got request for %d\nResults:\n", int32(*request.Point))

    mutex.RLock()
    matches := Data.Get(common.IntInterval{Start: int(*request.Point), End: int(*request.Point)})
    mutex.RUnlock()

	var response = messages.IntervalResponse{}
    for _, match := range matches {
		response.Points = append(response.Points, int32(match.ID()))
		fmt.Printf("Id: %d Range: %d - %d\n", match.ID(), match.Range().Start, match.Range().End)
    }
	fmt.Println()

	err = common.WriteMessage(conn, &response)
	if err != nil {
		log.Print(err)
		return
	}
}
