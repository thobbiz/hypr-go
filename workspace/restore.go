package workspace

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

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
		found := make(chan string, 1)
		done := make(chan struct{})

		go ipc.ListenEvents(done, func(event, data string) {
			if event == "openwindow" {
				parts := strings.SplitN(data, ",", 4)
				class := parts[2]
				if class == w.Class {
					found <- parts[0]
					close(done)
				}
			}
		})

		if workspaceID != nil {
			w.Workspace = *workspaceID
		}

		select {
		case address := <-found:
			log.Printf("address = %s----name = %s", address, w.Class)

			if w.Floating {
				err := setFloatingToTrue(address)
				if err != nil {
					log.Fatal(err)
				}

				err = dispatchWindow(address, w.Workspace)
				if err != nil {
					log.Fatal(err)
				}

				err = resizeWindow(address, w.Width, w.Height)
				if err != nil {
					log.Fatal(err)
				}

				err = moveWindow(address, w.X, w.Y)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				err := dispatchWindow(address, w.Workspace)
				if err != nil {
					log.Fatal(err)
				}
			}
		case <-time.After(10 * time.Second):
			log.Printf("timeout waiting for %s, skipping", w.Class)
			close(done)
			continue
		}
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

func dispatchWindow(address string, workspaceID int) error {
	_, err := ipc.SendCommand(fmt.Sprintf("dispatch hl.dsp.window.move({window = 'address:0x%s', workspace = '%d'})", address, workspaceID))
	if err != nil {
		return err
	}
	return nil
}

func setFloatingToTrue(address string) error {
	_, err := ipc.SendCommand(fmt.Sprintf("dispatch hl.dsp.window.float({window = 'address:0x%s', action = 'on'})", address))
	if err != nil {
		return err
	}
	return nil
}

func resizeWindow(address string, width, height int) error {
	_, err := ipc.SendCommand(fmt.Sprintf("dispatch hl.dsp.window.resize({window = 'address:0x%s', x = %d, y = %d})", address, width, height))
	if err != nil {
		return err
	}
	return nil
}

func moveWindow(address string, x, y int) error {
	_, err := ipc.SendCommand(fmt.Sprintf("dispatch hl.dsp.window.move({window = 'address:0x%s', x = %d, y = %d})", address, x, y))
	if err != nil {
		return err
	}
	return nil
}
