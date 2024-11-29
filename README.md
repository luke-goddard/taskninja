[![Go Report Card](https://goreportcard.com/badge/github.com/luke-goddard/taskninja)](https://goreportcard.com/report/github.com/luke-goddard/taskninja)

# TaskNinja

WIP TaskNinja is a command line task list utity program, inspired by
[taskwarrior](https://github.com/GothenburgBitFactory/taskwarrior)

![Screenshot](assets/screenshot.png?raw=true "Terminal User Interface Screenshot")

## Install

### From Git

```bash
git clone git@github.com:luke-goddard/taskninja.git
cd taskninja
make build
make install
```

## TUI Shortcuts

### TUI Shortcuts - Task Table

| Key | Action |
| --- | ------ |
| `q` | Quit |
| `a` | Add a new task |
| `d` | Complete a task |
| `e` | Edit a task |
| `f` | Filter tasks |
| `r` | Refresh the task list |
| `s` | Start The current task |
| `+` | Increase Priority|
| `-` | Decrease Priority|
| `H` | Set Priority to HIGH|
| `M` | Set Priority to MED|
| `L` | Set Priority to LOW|
| `N` | Set Priority to NONE|
| `Shift+D` | Delete a task |
| `g` | Go to the top row|
| `G` | Go to the bottom row|
| `/` | Fuzzy search |

## Configuration

Once TaskNinja has been installed, the first time you run the program it will
create a configuration file in your home directory. This file is located at
`~/.config/taskninja/config.yaml`. You can edit this file to change the default
settings.

### Sample Config

```yaml
connection:
    # File (sqlite), memory (sqlite)
    mode: file

    # Path to the database
    path: "/home/taskninaja/Documents/taskninja.db"

    # Default = path + ".bk"
    backupPath: "/home/taskninaja/Documents/taskninja.db.bk"

log:
    # debug, info, warn, error
    level: debug

    # Pretty, Json
    mode: json

    # Path to the log file: Default = /tmp/taskninja.log
    path: "/home/taskninja/Documents/taskninja.log"
```

## Local Development

### Install Optional Development Tools
```bash
go install github.com/air-verse/air@latest
go install github.com/onsi/ginkgo/v2/ginkgo
sudo add-apt-repository -y ppa:linuxgndu/sqlitebrowser-testing
sudo apt-get update
sudo apt-get install sqlitebrowser
```

### Run
```bash make run ```

### Run the unit tests
```bash make tests ```

### View the database
Note that pressing `ctrl+r` will refresh the database
```bash make browse ```

