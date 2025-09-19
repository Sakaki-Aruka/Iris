# Iris
Iris is a cli tool what manages minecraft server sessions.

# Commands
Base: `iris`

## Profile
syntax: 
- `iris profile [create | delete | update] [@filename]`

## Profile file format
- name: [This profile name]
- keep-logs: <true | false>
- startup-script: [script | @filepath]
- auto-restart
  - max: [number >= 0]
  - trigger-code: [exit code array (comma separated)]
  - startup-script: [script | @filepath]
  - scheduled-script
    - timing: [crontab date style]
    - profile: [profile name]

## Session
syntax: `iris session [list | create | connect | send | restart] [options...]`

### Options
- ((create) `--log-file=[filename]`)
- ((create) `--keep-logs`)
- (create) `--startup-script=[script | @filename]`
- (create) `--profile=[profile name]`
- (create) `--save-session`
- (create) `--auto-restart`
- (create) `--auto-restart-max=[number > 0]` (default = `3`)
- (create) `--auto-restart-trigger-code=[exit code array (comma separated)]` (default = `[1]`)
- (create) `--auto-restart-startup-script=[script | @filename]` (default = None)
- ((create) `--scheduled-script=[crontab date style] [profile]`)
- (connect | send | restart) `--sesssion=[session name]`
- (connect) `--previous=[lines >= 0 | all]` (default = `0`)
- (send) `--command=[command | @filename]`