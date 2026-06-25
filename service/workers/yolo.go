package worker

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/Voltamon/logr/data"
)

func StartYoloWorker(ctx context.Context, cameraID string, port string, logChan chan<- data.LogMessage) error {
    addr, err := net.ResolveUDPAddr("udp", port)
    if err != nil {
        return fmt.Errorf("Failed to resolve UDP address %s", err.Error())
    }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        return fmt.Errorf("Failed to listen on UDP port %s: %s", port, err.Error())
    }

    buffer := make([]byte, 2048)
    defer conn.Close()

    go func() {
        <-ctx.Done()
        conn.Close()
    }()

    for {
        nByte, _, err := conn.ReadFromUDP(buffer)
        if err != nil {
            return fmt.Errorf("Failed to read from UDP: %s", err.Error())
        }

        rawPayload := strings.TrimSpace(string(buffer[:nByte]))
        logChan <- data.YoloLogMessage(cameraID, data.LevelInfo, rawPayload)
        log.Printf("%s", rawPayload)
    }
}
