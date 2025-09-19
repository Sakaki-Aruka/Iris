package profile

type Profile struct {
	Name     string
	KeepLogs bool
}

type ScriptRef struct {
	Inline   string
	FilePath string
}

type AutoRestartCfg struct {
	Max           uint
	TriggerCodes  []int
	StartupScript ScriptRef
	Scheduled     []Schedule
}

type Schedule struct {
	Timing  string
	Profile string
}
