package ipc

type Monitor struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Height           int       `json:"height"`
	Width            int       `json:"width"`
	RefreshRate      float64   `json:"refreshRate"`
	ActiveWorkspace  Workspace `json:"activeWorkspace"`
	SpecialWorkspace Workspace `json:"specialWorkspace"`
	Scale            float64   `json:"scale"`
	Focused          bool      `json:"focused"`
	DpmsStatus       bool      `json:"dpmsStatus"`
}
