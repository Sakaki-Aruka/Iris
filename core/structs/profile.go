package structs

type Profile struct {
	Name          string      `json:"name"`
	KeepLogs      bool        `json:"keep_logs"`
	StartupScript ScriptOr    `json:"startup_script"`
	AutoRestart   AutoRestart `json:"auto_restart"`
}

type AutoRestart struct {
	Max           uint     `json:"max"`
	TriggerCode   []int    `json:"trigger_code"`
	StartupScript ScriptOr `json:"startup_script"`
	Schedule      Schedule `json:"schedule"`
}

type ScriptOr struct {
	Content string `json:"content"`
	IsFile  bool   `json:"is_file"`
}

type Schedule struct {
	Timing  string `json:"timing"`
	Profile string `json:"profile"`
}
