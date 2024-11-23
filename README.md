[![Go Report Card](https://goreportcard.com/badge/github.com/luke-goddard/taskninja)](https://goreportcard.com/report/github.com/luke-goddard/taskninja)

# TaskNinja

WIP TaskNinja is a command line task list utity program, inspired by
[taskwarrior](https://github.com/GothenburgBitFactory/taskwarrior)

## Install

### From Git

```bash
git clone git@github.com:luke-goddard/taskninja.git
cd taskninja
make build
make install
```

## Configuration

Once TaskNinja has been installed, the first time you run the program it will
create a configuration file in your home directory. This file is located at
`~/.config/taskninja/config.yaml`. You can edit this file to change the default
settings.


## TUI Shortcuts

## TUI Shortcuts - Task Table

| Key | Action |
| --- | ------ |
| `q` | Quit |
| `a` | Add a new task |
| `d` | Delete a task |
| `e` | Edit a task |
| `f` | Filter tasks |
| `r` | Refresh the task list |
| `s` | Start The current task |

## Local Development

### Install Optional Development Tools
```bash
go install github.com/air-verse/air@latest
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
