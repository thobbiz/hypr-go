package workspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/thobbiz/hypr-go/appindex"
	"github.com/thobbiz/hypr-go/hyprerr"
	"github.com/thobbiz/hypr-go/ipc"
	"gopkg.in/yaml.v3"
)

func SaveWorkspace(id *int, name string, dir string) error {
	path, err := LayoutPath(dir, name)
	if err != nil {
		return err
	}

	var wsID int
	if id == nil {
		current, err := getCurrentWorkspaceId()
		if err != nil {
			return err
		}
		wsID = *current
	} else {
		wsID = *id
	}

	windows, err := getWorkspaceWindows(wsID)
	if err != nil {
		return err
	}

	layout := WorkspaceLayout{
		Name:    name,
		SavedAt: time.Now(),
		Windows: *windows,
	}

	return saveLayoutAsYaml(layout, path)

}

func saveLayoutAsYaml(layout WorkspaceLayout, path string) error {
	data, err := yaml.Marshal(layout)
	if err != nil {
		return hyprerr.New("saveLayoutAsYaml", hyprerr.ErrParse, err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return hyprerr.New("saveLayoutAsYaml", hyprerr.ErrIO, err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
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
		return nil, hyprerr.New("unmarshalJSON", hyprerr.ErrParse, err)
	}

	var onWorkspace []ipc.Window
	seen := map[string]struct{}{}
	var classes []string
	addClass := func(c string) {
		if c == "" {
			return
		}
		if _, ok := seen[c]; !ok {
			seen[c] = struct{}{}
			classes = append(classes, c)
		}
	}

	for _, w := range windows {
		if w.Workspace.ID != id {
			continue
		}
		onWorkspace = append(onWorkspace, w)
		addClass(w.Class)
		addClass(w.InitialClass)
	}

	if len(onWorkspace) == 0 {
		return nil, hyprerr.New("getWorkspaceWindows", hyprerr.ErrNotFound, fmt.Errorf("no windows found on workspace %d", id))
	}

	execs, err := appindex.ResolveExecs(classes)
	if err != nil {
		// Not fatal: every window still gets the pid fallback below.
		log.Printf("desktop file resolution failed: %v", err)
	}

	var result []WindowEntry
	for _, w := range onWorkspace {
		exec := execs[w.Class]
		if exec == "" {
			exec = execs[w.InitialClass]
		}

		if exec == "" {
			log.Printf("%s: no desktop entry, falling back to pid %d", w.Class, w.PID)
			exec, err = getCommandFromPID(w.PID)
			if err != nil {
				log.Printf("skipping %s: %v", w.Class, err)
				continue
			}
		}

		result = append(result, WindowEntry{
			Class:     w.Class,
			Exec:      exec,
			Workspace: w.Workspace.ID,
			X:         w.At[0],
			Y:         w.At[1],
			Width:     w.Size[0],
			Height:    w.Size[1],
			Floating:  w.Floating,
		})
	}

	if len(result) == 0 {
		return nil, hyprerr.New("getWorkspaceWindows", hyprerr.ErrNotFound, fmt.Errorf("no launchable windows on workspace %d", id))
	}

	return &result, nil
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

func LayoutPath(dir, name string) (string, error) {
	if err := validateName(name); err != nil {
		return "", err
	}

	base, err := resolveDir(dir)
	if err != nil {
		return "", err
	}

	return filepath.Join(base, name+".yaml"), nil
}

func validateName(name string) error {
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		log.Println(name)
		return hyprerr.New("validateName", hyprerr.ErrInvalidInput, fmt.Errorf("invalid layout name %q", name))
	}
	return nil
}

func resolveDir(dir string) (string, error) {
	if dir == "" {
		return defaultDir()
	}

	if dir == "~" || strings.HasPrefix(dir, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", hyprerr.New("resolveDir", hyprerr.ErrIO, err)
		}
		dir = filepath.Join(home, strings.TrimPrefix(dir, "~"))
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", hyprerr.New("resolveDir", hyprerr.ErrIO, err)
	}
	return abs, nil
}

func defaultDir() (string, error) {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", hyprerr.New("defaultDir", hyprerr.ErrIO, err)
		}
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(base, "hypr-go", "workspace"), nil
}
