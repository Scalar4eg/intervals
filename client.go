package main

import (
	"net"
	"github.com/bo0rsh201/intervals/proto"
	"github.com/bo0rsh201/intervals/common"
	"log"
)

func TestClient(n int32) {
	conn, _ := net.Dial("tcp", "127.0.0.1:8080")
	for {
		request := messages.IntervalRequest{}
		request.Point = &n

		err := common.WriteMessage(conn, &request)
		if err != nil {
			log.Print(err)
			return
		}
		response := messages.IntervalResponse{}
		err = common.ReadMessage(conn, &response)
		if err != nil {
			log.Print(err)
			return
		}
		log.Print(response.String())
	}

}
