package structs

type ConnectionRequest struct {
	SessionName string   `json:"session_name"`
	Env         []string `json:"env"`
	Pwd         string   `json:"pwd"`
	Command     string   `json:"command"`
}
