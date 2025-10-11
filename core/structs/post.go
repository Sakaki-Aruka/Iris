package structs

const (
	UserPost = iota
	ProcessPost
)

type Post struct {
	Type        int    `json:"type"`
	SessionName string `json:"session_name"`
	Line        []byte `json:"line"`
}
