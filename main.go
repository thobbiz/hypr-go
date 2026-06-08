package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/thobbiz/hypr-go/ipc"
	"github.com/thobbiz/hypr-go/layout"
)

func main() {
	err := layout.SaveWorkspace(1, "coded")
	if err != nil {
		log.Fatal(err)
	}
}

func get_clients() ([]ipc.Window, error) {
	out, err := exec.Command("hyprctl", "-j", "clients").Output()
	if err != nil {
		return nil, err
	}

	var windows []ipc.Window
	if err := json.Unmarshal(out, &windows); err != nil {
		return nil, err
	}

	return windows, nil
}

func print_events() {
	his := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	xdg := os.Getenv("XDG_RUNTIME_DIR")

	path := fmt.Sprintf("%s/hypr/%s/.socket2.sock", xdg, his)

	conn, err := net.Dial("unix", path)
	if err != nil {

	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()

		fmt.Println(line)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from socket: %v", err)
	}
}

func change_workspace() {
	his := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	xdg := os.Getenv("XDG_RUNTIME_DIR")

	path := fmt.Sprintf("%s/hypr/%s/.socket.sock", xdg, his)

	conn, err := net.Dial("unix", path)
	if err != nil {
		// handle error
	}
	fmt.Fprintf(conn, "/dispatch workspace 4")
	status, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(status)
}
