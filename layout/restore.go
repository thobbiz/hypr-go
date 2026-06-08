package layout

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/thobbiz/hypr-go/ipc"
	"gopkg.in/yaml.v3"
)

func RestoreWorkspace(name string, workspaceID *int) error {
	if workspaceID != nil && *workspaceID <= 0 {
		return errors.New("workspace id cannot be zero(0) or negative")
	}

	saved, err := decodeYaml(name)
	if err != nil {
		log.Fatal(err)
	}

	for _, w := range saved.Windows {
		LaunchAppFromCmd(w)
		found := make(chan string)

		go ipc.ListenEvents(func(event, data string) {
			if event == "openwindow" {
				parts := strings.SplitN(data, ",", 4)
				class := parts[2]
				if class == w.Class {
					found <- parts[0]
				}
			}
		})

		address := <-found

		if workspaceID != nil {
			w.Workspace = *workspaceID
		}
		out, err := dispatchWindow(address, w.Workspace)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(out)
	}

	return nil
}

func LaunchAppFromCmd(entry WindowEntry) {
	parts := strings.Fields(entry.Cmd)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Start()
}

func decodeYaml(filename string) (*Layout, error) {
	name := fmt.Sprintf("./data/%s.yaml", filename)
	f, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var layout Layout

	if err := yaml.Unmarshal(f, &layout); err != nil {
		return nil, err
	}

	return &layout, nil
}

func dispatchWindow(address string, workspaceID int) (string, error) {
	out, err := ipc.SendCommand(fmt.Sprintf("dispatch hl.dsp.window.move({window = 'address:0x%s', workspace = '%d'})", address, workspaceID))
	if err != nil {
		return "", err
	}
	return out, nil
}
