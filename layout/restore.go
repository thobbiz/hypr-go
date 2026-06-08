package layout

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

func RestoreWorkspace(name string, workspaceID int) {

}

func LaunchAppFromCmd(entry WindowEntry) {
	parts := strings.Fields(entry.Cmd)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Start()
	log.Println("success jor")
}

func DecodeYaml(filename string) (*Layout, error) {
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

func Dispatch() {

}
