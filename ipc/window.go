package ipc

type Window struct {
	Address        string          `json:"address"`
	Class          string          `json:"class"`
	Mapped         bool            `json:"mapped"`
	Hidden         bool            `json:"hidden"`
	At             [2]int          `json:"at"`
	Size           [2]int          `json:"size"`
	Workspace      WindowWorkspace `json:"workspace"`
	Floating       bool            `json:"floating"`
	Monitor        int             `json:"monitor"`
	Title          string          `json:"title"`
	PID            int             `json:"pid"`
	FullScreen     int             `json:"fullscreen"`
	FocusHistoryID int             `json:"focus_history_id"`
	ContentType    string          `json:"content_type"`
	StableID       int             `json:"stable_id"`
}

type WindowWorkspace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
