package workspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/thobbiz/hypr-go/hyprerr"
	"github.com/thobbiz/hypr-go/ipc"
	"gopkg.in/yaml.v3"
)

func SaveWorkspace(id *int, name string) error {
	var ID *int
	var err error
	if id == nil {
		ID, err = getCurrentWorkspaceId()
		if err != nil {
			return err
		}
	} else {
		ID = id
	}

	windows, err := getWorkspaceWindows(*ID)
	if err != nil {
		return err
	}

	layout := Layout{
		Name:    name,
		SavedAt: time.Now(),
		Windows: *windows,
	}

	err = saveLayoutAsYaml(layout)
	if err != nil {
		return err
	}

	return nil
}

func saveLayoutAsYaml(layout Layout) error {
	data, err := yaml.Marshal(layout)
	if err != nil {
		return hyprerr.New("saveLayoutAsYaml", hyprerr.ErrParse, err)
	}

	name := fmt.Sprintf("./data/workspace/%s.yaml", layout.Name)
	err = os.WriteFile(name, data, 0644)
	if err != nil {
		return hyprerr.New("saveLayoutAsYaml", hyprerr.ErrIO, err)
	}

	return nil
}

func getWorkspaceWindows(id int) (*[]WindowEntry, error) {
	resp, err := ipc.SendCommand("j/clients")
	if err != nil {
		return nil, hyprerr.New("getWorkspaceWindows", hyprerr.ErrIPC, err)
	}

	var windows []ipc.Window
	if err := json.Unmarshal([]byte(resp), &windows); err != nil {
		return nil, hyprerr.New("getWorkspaceWindows", hyprerr.ErrParse, err)
	}

	var result []WindowEntry

	for _, w := range windows {
		if w.Workspace.ID == id {
			command, err := getCommandFromPID(w.PID)
			if err != nil {
				return nil, err
			}
			anon := WindowEntry{
				Class:     w.Class,
				Cmd:       command,
				Workspace: w.Workspace.ID,
				X:         w.At[0],
				Y:         w.At[1],
				Width:     w.Size[0],
				Height:    w.Size[1],
				Floating:  w.Floating,
			}
			result = append(result, anon)
		}
	}

	if len(result) == 0 {
		return nil, hyprerr.New("getWorkspaceWindows", hyprerr.ErrNotFound, fmt.Errorf("no windows found on workspace %d", id))
	}
	return &result, nil
}

func getCommandFromPID(pid int) (string, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return "", hyprerr.New("getCommandFromPID", hyprerr.ErrIO, err)
	}

	parts := bytes.Split(data, []byte{0})
	if len(parts) == 0 || len(parts[0]) == 0 {
		return "", hyprerr.New("getCommandFromPID", hyprerr.ErrParse, fmt.Errorf("empty cmdline for pid %d", pid))
	}

	exe, _, _ := strings.Cut(string(parts[0]), " ")
	return exe, nil
}

func getCurrentWorkspaceId() (*int, error) {
	raw, err := ipc.SendCommand("j/activeworkspace")
	if err != nil {
		return nil, hyprerr.New("getCurrentWorkspaceId", hyprerr.ErrIPC, err)
	}
	var ws ipc.Workspace
	if err := json.Unmarshal([]byte(raw), &ws); err != nil {
		return nil, hyprerr.New("getCurrentWorkspaceId", hyprerr.ErrParse, err)
	}
	return &ws.ID, nil
}
