package ipc

type Monitor struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Height           int       `json:"height"`
	Width            int       `json:"width"`
	RefreshRate      float64   `json:"refreshRate"`
	X                int       `json:"x"`
	Y                int       `json:"y"`
	Scale            float64   `json:"scale"`
	Transform        int       `yaml:"transform"`
	Focused          bool      `json:"focused"`
	ActiveWorkspace  Workspace `json:"activeWorkspace"`
	SpecialWorkspace Workspace `json:"specialWorkspace"`
	DpmsStatus       bool      `json:"dpmsStatus"`
}
