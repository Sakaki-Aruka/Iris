# Iris
Iris is a cli tool what manages minecraft server sessions.

# Commands
Base: `iris`

## Profile
syntax: 
- `iris profile [create | delete | update] [@filename]`
- `iris profile list`

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
- (create) `--name=[name]`
- (create) `--startup-script=[script]`
- (create) `--profile=[profile name]`
- (create) `--save-session`
- (create) `--auto-restart`
- (create) `--auto-restart-max=[number > 0]` (default = `3`)
- (create) `--auto-restart-trigger-code=[exit code array (comma separated)]` (default = `[]`)
- (create) `--auto-restart-startup-script=[script]` (default = None)
- ((create) `--scheduled-script-timing=[crontab date style]`)
- ((create) `--scheduled-script=[script]`)
- (connect | send | restart) `--sesssion=[session name]`
- (connect) `--previous=[lines >= 0 | all]` (default = `0`)
- (send) `--command=[command | @filename]`

## Template
syntax: `iris template [profile | systemd]`

## Sys
syntax: `iris daemon [start | (stop)]`