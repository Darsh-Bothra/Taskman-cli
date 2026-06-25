# taskman

A fast, local CLI task manager written in Go.

Tasks are stored in a plain JSON file — no database, no internet, no fuss.

---

## Installation

```bash
go install github.com/darshbothra007/taskman@latest
```

Or build from source:

```bash
git clone https://github.com/darshbothra007/taskman
cd taskman
make install
```

---

## Usage

```
taskman [command] [flags]
```

### Commands

| Command | Description |
|---|---|
| `add <description>` | Add a new task |
| `list` | List pending tasks |
| `list --all` | List all tasks including completed |
| `get <id>` | Show full details of a task |
| `complete <id>` | Mark a task as done |
| `delete <id>` | Remove a task permanently |

### Examples

```bash
# Add tasks
taskman add "Write unit tests"
taskman add Buy groceries and cook dinner

# List pending tasks
taskman list

# List everything
taskman list --all

# Inspect one task
taskman get 3

# Mark as done
taskman complete 3

# Remove a task
taskman delete 5
```

### Global flags

| Flag | Default | Description |
|---|---|---|
| `--file` | `~/.taskman.json` | Path to the task storage file |

---

## Configuration

taskman reads from `~/.config/taskman/taskman.yaml` (optional).

```yaml
# ~/.config/taskman/taskman.yaml
file: /path/to/custom/tasks.json
```

Environment variables are also supported with the `TASKMAN_` prefix:

```bash
TASKMAN_FILE=/tmp/work-tasks.json taskman list
```

---

## Project structure

```
taskman/
├── main.go                   # Entry point
├── cmd/                      # Cobra commands (one file per command)
│   ├── root.go               # Root command, lifecycle hooks
│   ├── add.go
│   ├── list.go
│   ├── get.go
│   ├── complete.go
│   └── delete.go
├── internal/
│   ├── task/                 # Domain logic — Task type & Manager
│   │   ├── task.go
│   │   └── task_test.go
│   ├── storage/              # JSON persistence
│   │   ├── storage.go
│   │   └── storage_test.go
│   └── config/               # Viper configuration
│       └── config.go
├── Makefile
└── go.mod
```

---

## Development

```bash
make build        # compile to ./bin/taskman
make test         # run tests with race detector
make test-cover   # generate HTML coverage report
make lint         # run golangci-lint
make clean        # remove build artefacts
```

---

## License

MIT © Darsh Bothra