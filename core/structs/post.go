package structs

const (
	User = iota
	Process
)

type Post struct {
	Type        int    `json:"type"`
	SessionName string `json:"session_name"`
	Line        []byte `json:"line"`
}
