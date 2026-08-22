package ipc

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

func SendCommand(command string) (string, error) {
	his := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	xdg := os.Getenv("XDG_RUNTIME_DIR")

	path := fmt.Sprintf("%s/hypr/%s/.socket.sock", xdg, his)

	conn, err := net.Dial("unix", path)
	if err != nil {
		return "", fmt.Errorf("failed to create connection: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(command))
	if err != nil {
		return "", fmt.Errorf("failed to write to connection: %v", err)
	}

	data, err := io.ReadAll(conn)
	if err != nil {
		return string(data), fmt.Errorf("failed to read from connection: %v", err)
	}

	return string(data), nil
}

func ListenEvents(ctx context.Context, handler func(event, data string)) error {
	his := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	xdg := os.Getenv("XDG_RUNTIME_DIR")
	path := fmt.Sprintf("%s/hypr/%s/.socket2.sock", xdg, his)
	conn, err := net.Dial("unix", path)
	if err != nil {
		return fmt.Errorf("failed to create connection: %v", err)
	}
	defer conn.Close()

	go func() {
		<-ctx.Done()
		conn.SetReadDeadline(time.Now())
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ">>", 2)
		if len(parts) == 2 {
			handler(parts[0], parts[1])
		}
	}
	if ctx.Err() != nil {
		return nil
	}
	return scanner.Err()
}
