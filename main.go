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
)

func main() {
	// windows, err := get_clients()
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	return
	// }

	// for _, w := range windows {
	// 	fmt.Printf("Title: %s | Mapped: %t | Workspace: %v | PID: %d\n", w.Title, w.Mapped, w.Workspace, w.PID)
	// }

	err := ipc.ListenEvents(func(event, data string) {
		fmt.Printf("Event: %s | Data: %s\n", event, data)
	})
	if err != nil {
		log.Fatal(err)
	}
}

func get_clients() ([]Window, error) {
	out, err := exec.Command("hyprctl", "-j", "clients").Output()
	if err != nil {
		return nil, err
	}

	var windows []Window
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
