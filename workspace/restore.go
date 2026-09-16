package workspace

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/thobbiz/hypr-go/hyprerr"
	"github.com/thobbiz/hypr-go/ipc"
	"gopkg.in/yaml.v3"
)

func RestoreWorkspace(name string, workspaceID *int) error {
	if workspaceID != nil && *workspaceID <= 0 {
		return hyprerr.New("RestoreWorkspace", hyprerr.ErrInvalidInput, fmt.Errorf("workspace id cannot be zero(0) or negative"))
	}
	saved, err := decodeYaml(name)
	if err != nil {
		return err
	}

	for _, w := range saved.Windows {
		launchAppFromExec(w.Exec)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		found := make(chan string, 1)

		go func(targetClass string) {
			_ = ipc.ListenEvents(ctx, func(event, data string) {
				if event == "openwindow" {
					parts := strings.SplitN(data, ",", 4)
					if len(parts) >= 3 && parts[2] == targetClass {
						found <- parts[0]
						cancel()
					}
				}
			})
		}(w.Class)

		if workspaceID != nil {
			w.Workspace = *workspaceID
		}

		select {
		case address := <-found:
			cancel()
			log.Printf("address = %s----name = %s", address, w.Class)

			if w.Floating {
				err := setFloatingToTrue(address)
				if err != nil {
					log.Println(err)
				}
				err = dispatchWindow(address, w.Workspace)
				if err != nil {
					log.Println(err)
				}
				err = resizeWindow(address, w.Width, w.Height)
				if err != nil {
					log.Println(err)
				}
				err = moveWindow(address, w.X, w.Y)
				if err != nil {
					log.Println(err)
				}
			} else {
				err := dispatchWindow(address, w.Workspace)
				if err != nil {
					log.Println(err)
				}
			}
		case <-ctx.Done():
			cancel()
			log.Printf("timeout waiting for %s, skipping", w.Class)
			continue
		}
	}
	return nil
}

func launchAppFromExec(cmdStr string) {
	parts := strings.Fields(cmdStr)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Start()
}

func decodeYaml(filename string) (*Layout, error) {
	name := fmt.Sprintf("./data/workspace/%s.yaml", filename)
	f, err := os.ReadFile(name)
	if err != nil {
		return nil, hyprerr.New("decodeYaml", hyprerr.ErrIO, err)
	}

	var layout Layout

	if err := yaml.Unmarshal(f, &layout); err != nil {
		return nil, hyprerr.New("decodeYaml", hyprerr.ErrParse, err)
	}

	return &layout, nil
}

func dispatchWindow(address string, workspaceID int) error {
	_, err := ipc.SendCommand(fmt.Sprintf("dispatch hl.dsp.window.move({window = 'address:0x%s', workspace = '%d'})", address, workspaceID))
	if err != nil {
		return hyprerr.New("dispatchWindow", hyprerr.ErrIPC, err)
	}
	return nil
}

func setFloatingToTrue(address string) error {
	_, err := ipc.SendCommand(fmt.Sprintf("dispatch hl.dsp.window.float({window = 'address:0x%s', action = 'on'})", address))
	if err != nil {
		return hyprerr.New("setFloatingToTrue", hyprerr.ErrIPC, err)
	}
	return nil
}

func resizeWindow(address string, width, height int) error {
	_, err := ipc.SendCommand(fmt.Sprintf("dispatch hl.dsp.window.resize({window = 'address:0x%s', x = %d, y = %d})", address, width, height))
	if err != nil {
		return hyprerr.New("resizeWindow", hyprerr.ErrIPC, err)
	}
	return nil
}

func moveWindow(address string, x, y int) error {
	_, err := ipc.SendCommand(fmt.Sprintf("dispatch hl.dsp.window.move({window = 'address:0x%s', x = %d, y = %d})", address, x, y))
	if err != nil {
		return hyprerr.New("moveWindow", hyprerr.ErrIPC, err)
	}
	return nil
}
