package ipc

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func SendCommand(command string) (string, error) {
	his := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	xdg := os.Getenv("XDG_RUNTIME_DIR")

	path := fmt.Sprintf("%s/hypr/%s/.socket.sock", xdg, his)

	conn, err := net.Dial("unix", path)
	if err != nil {
		return "", fmt.Errorf("Error creating connection: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(command))
	if err != nil {
		return "", fmt.Errorf("Error writing command connection: %v", err)
	}

	data, err := io.ReadAll(conn)
	if err != nil {
		return string(data), fmt.Errorf("Error reading from socket: %v", err)
	}

	return string(data), nil
}

func ListenEvents(handler func(event, data string)) error {
	his := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	xdg := os.Getenv("XDG_RUNTIME_DIR")

	path := fmt.Sprintf("%s/hypr/%s/.socket2.sock", xdg, his)

	conn, err := net.Dial("unix", path)
	if err != nil {
		return fmt.Errorf("Error creating connection: %v", err)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ">>")
		if len(parts) == 2 {
			handler(parts[0], parts[1])
		}
	}

	return scanner.Err()
}
