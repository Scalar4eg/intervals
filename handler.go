package main
import (
    "net"
    "log"
    "bufio"
    "strconv"
    "strings"
    "encoding/json"
    "github.com/bo0rsh201/intervals/common"
    "bytes"
    "fmt"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)
    req_str, err := reader.ReadString('\n')
    if err != nil {
        log.Print(err)
        return
    }

    dst, err := strconv.Atoi(strings.TrimRight(req_str, "\n\r"))
    if err != nil {
        log.Print(err)
        return
    }
    var debugBuffer bytes.Buffer
    debugBuffer.WriteString(fmt.Sprintf("Got request for %d\nResults:\n", dst))
    mutex.RLock()
    matches := Data.Get(common.IntInterval{Start: dst, End: dst})
    mutex.RUnlock()
    result := make([]uintptr, len(matches))
    i := 0
    for _, match := range matches {
        debugBuffer.WriteString(fmt.Sprintf("Id: %d Range: %d - %d\n", match.ID(), match.Range().Start, match.Range().End))
        result[i] = match.ID()
        i++
    }
    debugBuffer.WriteString("\n")
    json_bytes, err := json.Marshal(result)
    if err != nil {
        log.Print(err)
        return
    }

    _, err = conn.Write(append(json_bytes, '\n'))
    if err != nil {
        log.Print(err)
        return
    }
    fmt.Print(debugBuffer.String())
}
