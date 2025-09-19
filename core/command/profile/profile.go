package profile

type Profile struct {
	Name     string `json:"name"`
	KeepLogs bool   `json:"keep_logs"`
}

type ScriptRef struct {
	Inline   string `json:"inline"`
	FilePath string `json:"file_path"`
}

type AutoRestartCfg struct {
	Max           uint       `json:"max"`
	TriggerCodes  []int      `json:"trigger_codes"`
	StartupScript ScriptRef  `json:"startup_script"`
	Scheduled     []Schedule `json:"scheduled"`
}

type Schedule struct {
	Timing  string `json:"timing"`
	Profile string `json:"profile"`
}
