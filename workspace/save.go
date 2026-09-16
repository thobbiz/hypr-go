package workspace

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/thobbiz/hypr-go/hyprerr"
	"github.com/thobbiz/hypr-go/ipc"
	"gopkg.in/yaml.v3"
)

var XDGFIELDCODES = regexp.MustCompile(`%[fFuUdDnNickvm]`)

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
		return nil, hyprerr.New("unmarshalJSON", hyprerr.ErrParse, err)
	}

	var result []WindowEntry

	for _, w := range windows {
		if w.Workspace.ID == id {
			exec, err := getExecCommand(w)
			if err != nil {
				return nil, hyprerr.New("getExecCommand", hyprerr.ErrNotFound, err)
			}
			anon := WindowEntry{
				Class:     w.Class,
				Exec:      exec,
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

func getExecCommand(w ipc.Window) (string, error) {
	dirs := getDesktopDirs()

	desktopFile := findDesktopFile(w.Class, dirs)
	if desktopFile == "" && w.InitialClass != "" {
		desktopFile = findDesktopFile(w.InitialClass, dirs)
	}

	execCmd := ""
	if desktopFile != "" {
		execCmd = extractExec(desktopFile)
	}

	if execCmd == "" {
		log.Printf("%v couldn't extract exec from desktop, switching to pid...", w.Class)
		execCmd, err := getCommandFromPID(w.PID)
		if err != nil {
			return "", fmt.Errorf("exec command not found for %v", w.Class)
		}
		return execCmd, nil
	}

	return execCmd, nil
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

func getDesktopDirs() []string {
	home, _ := os.UserHomeDir()
	return []string{
		filepath.Join(home, ".local/share/applications"),
		"/usr/share/applications",
		"/usr/local/share/applications",
		"/var/lib/flatpak/exports/share/applications",
		filepath.Join(home, ".local/share/flatpak/exports/share/applications"),
	}
}

// extractExec parses the Exec line from the [Desktop Entry] section
func extractExec(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inDesktopEntry := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			inDesktopEntry = (line == "[Desktop Entry]")
			continue
		}

		if inDesktopEntry && strings.HasPrefix(line, "Exec=") {
			cmd := strings.TrimPrefix(line, "Exec=")
			// Remove XDG field codes (%u, %F, etc.)
			cmd = XDGFIELDCODES.ReplaceAllString(cmd, "")
			return strings.TrimSpace(cmd)
		}
	}
	return ""
}

// findDesktopFile attempts to match class against StartupWMClass or filename
func findDesktopFile(targetClass string, dirs []string) string {
	if targetClass == "" {
		return ""
	}

	var fallbackMatch string

	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range files {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".desktop") {
				continue
			}

			fullPath := filepath.Join(dir, entry.Name())
			fileBytes, err := os.ReadFile(fullPath)
			if err != nil {
				continue
			}
			content := string(fileBytes)

			// Direct match on StartupWMClass
			startupWMRegex := regexp.MustCompile(`(?m)^StartupWMClass=(?i:` + regexp.QuoteMeta(targetClass) + `)$`)
			if startupWMRegex.MatchString(content) {
				return fullPath
			}

			// Direct match on filename
			base := strings.TrimSuffix(entry.Name(), ".desktop")
			if strings.EqualFold(base, targetClass) {
				return fullPath
			}

			// Fallback partial match if nothing exact found
			if fallbackMatch == "" && strings.Contains(strings.ToLower(base), strings.ToLower(targetClass)) {
				fallbackMatch = fullPath
			}
		}
	}

	return fallbackMatch
}
