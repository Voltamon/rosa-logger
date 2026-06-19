package worker

import (
	"context"
	"fmt"
	"io"

	"github.com/Voltamon/logr/data"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

type ChannelWriter struct {
	containerID string
	level       data.LogLevel
	logChan     chan<- data.LogMessage
}

func (cw *ChannelWriter) Write(p []byte) (n int, err error) {
	logPayload := string(p)
	cw.logChan <- data.DockLogMessage(cw.containerID, cw.level, logPayload)
	return len(p), nil
}

func StartDockWorker(ctx context.Context, containerName []string, logStream chan<- data.LogMessage) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	logOpts := client.ContainerLogsOptions{
		Timestamps: false,
		Tail: "5",
		Follow:     true,
		ShowStdout: true,
		ShowStderr: true,
	}

	containerData, err := cli.ContainerInspect(ctx, containerName[0], client.ContainerInspectOptions{})
	if err != nil {
		return fmt.Errorf("failed to inspect container: %w", err)
	}
	containerID := containerData.Container.ID

	stream, err := cli.ContainerLogs(ctx, containerID, logOpts)
	if err != nil {
		return fmt.Errorf("failed to get container logs: %w", err)
	}
	defer stream.Close()

	stdoutWriter := &ChannelWriter{
	    containerID:    containerID,
		level:          data.LevelInfo,
		logChan:        logStream,
	}

	stderrWriter := &ChannelWriter{
	    containerID: containerID,
		level: data.LevelError,
		logChan: logStream,
	}

	_, err = stdcopy.StdCopy(stdoutWriter, stderrWriter, stream)
	if err != nil && err != io.EOF {
		return fmt.Errorf("docker log stream error: %w", err)
	}

	return nil
}
