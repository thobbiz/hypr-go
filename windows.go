package main

type Window struct {
	Mapped         bool      `json:"mapped"`
	Hidden         bool      `json:"hidden"`
	Workspace      Workspace `json:"workspace"`
	Floating       bool      `json:"floating"`
	Monitor        int       `json:"monitor"`
	Title          string    `json:"title"`
	PID            int       `json:"pid"`
	FullScreen     int       `json:"fullscreen"`
	FocusHistoryID int       `json:"focus_history_id"`
	ContentType    string    `json:"content_type"`
	StableID       int       `json:"stable_id"`
}
