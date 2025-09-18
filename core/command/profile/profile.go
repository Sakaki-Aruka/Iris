package profile

type Profile struct {
	name          string
	keepLogs      bool
	startupScript ScriptOr
	autoRestart   AutoRestart
}

type AutoRestart struct {
	max           uint
	triggerCode   []int
	startupScript ScriptOr
	schedule      Schedule
}

type ScriptOr struct {
	content string
	isFile  bool
}

type Schedule struct {
	timing  string
	profile string
}
