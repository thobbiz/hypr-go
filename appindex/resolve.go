package appindex

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/thobbiz/hypr-go/hyprerr"
)

// Matches XDG field codes (%u, %F, ...) and the flatpak file markers (@@, @@u, @@f).
var fieldCodes = regexp.MustCompile(`%[fFuUdDnNickvm]|@@[uf]?`)

type entry struct {
	id      string
	wmClass string
	exec    string
	binary  string
}

// ResolveExecs maps each window class to the Exec command of its .desktop file
func ResolveExecs(classes []string) (map[string]string, error) {
	index, err := loadIndex()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(classes))
	for _, class := range classes {
		if e, ok := match(index, class); ok {
			result[class] = e.exec
		}
	}
	return result, nil
}

// match tries the strictest rule first across the whole index, then falls
// back: StartupWMClass, then desktop id, then binary name.
func match(index []entry, class string) (entry, bool) {
	if class == "" {
		return entry{}, false
	}

	tiers := []func(entry) string{
		func(e entry) string { return e.wmClass },
		func(e entry) string { return e.id },
		func(e entry) string { return e.binary },
	}

	for _, field := range tiers {
		for _, e := range index {
			if strings.EqualFold(field(e), class) {
				return e, true
			}
		}
	}
	return entry{}, false
}

func desktopDirs() []string {
	home, _ := os.UserHomeDir()
	return []string{
		filepath.Join(home, ".local/share/applications"),
		"/usr/local/share/applications",
		"/usr/share/applications",
		filepath.Join(home, ".local/share/flatpak/exports/share/applications"),
		"/var/lib/flatpak/exports/share/applications",
	}
}

// loadIndex reads every .desktop file once. Directories are in priority
// order, and the first file with a given name wins, so user entries override
// system ones.
func loadIndex() ([]entry, error) {
	seen := map[string]bool{}
	var index []entry
	readable := false

	for _, dir := range desktopDirs() {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		readable = true

		for _, f := range files {
			name := f.Name()
			if f.IsDir() || !strings.HasSuffix(name, ".desktop") || seen[name] {
				continue
			}
			seen[name] = true

			if e, ok := parseDesktopFile(filepath.Join(dir, name)); ok {
				index = append(index, e)
			}
		}
	}

	if !readable {
		return nil, hyprerr.New("loadIndex", hyprerr.ErrIO, fmt.Errorf("no readable application directories"))
	}
	return index, nil
}

func parseDesktopFile(path string) (entry, bool) {
	file, err := os.Open(path)
	if err != nil {
		return entry{}, false
	}
	defer file.Close()

	e := entry{id: strings.TrimSuffix(filepath.Base(path), ".desktop")}
	inEntry := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "[") {
			if inEntry {
				break // left the [Desktop Entry] group
			}
			inEntry = line == "[Desktop Entry]"
			continue
		}
		if !inEntry {
			continue
		}

		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		switch key {
		case "Type":
			if val != "Application" {
				return entry{}, false
			}
		case "Exec":
			e.exec = stripFieldCodes(val)
		case "StartupWMClass":
			e.wmClass = val
		}
	}

	// Check if the scan loop aborted due to an I/O error or line overflow
	if err := scanner.Err(); err != nil {
		return entry{}, false
	}

	if e.exec == "" {
		return entry{}, false
	}
	e.binary = binaryName(e.exec)
	return e, true
}

func stripFieldCodes(s string) string {
	return strings.Join(strings.Fields(fieldCodes.ReplaceAllString(s, "")), " ")
}

// binaryName returns the base name of the executable, skipping a leading
// "env" and any VAR=value assignments.
func binaryName(exec string) string {
	for tok := range strings.FieldsSeq(exec) {
		if tok == "env" || strings.Contains(tok, "=") {
			continue
		}
		return filepath.Base(tok)
	}
	return ""
}
