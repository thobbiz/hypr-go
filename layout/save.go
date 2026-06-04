package layout

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/thobbiz/hypr-go/ipc"
	"gopkg.in/yaml.v3"
)

func SaveWorkspace(id int, name string) error {
	windows, err := getWorkspaceWindows(id)
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
		log.Fatal(err)
	}

	name := fmt.Sprintf("./data/%s.yaml", layout.Name)
	err = os.WriteFile(name, data, 0644)
	if err != nil {
		return fmt.Errorf("Error creating file: %v", err)
	}

	return nil
}

func getWorkspaceWindows(id int) (*[]WindowEntry, error) {
	resp, err := ipc.SendCommand("j/clients")
	if err != nil {
		log.Fatal(err)
	}

	var windows []ipc.Window
	if err := json.Unmarshal([]byte(resp), &windows); err != nil {
		return nil, err
	}

	var result []WindowEntry

	for _, w := range windows {
		if w.Workspace.ID == id {
			anon := WindowEntry{
				Class:     w.Class,
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
		return nil, errors.New("workspace not found")
	}
	return &result, nil
}
