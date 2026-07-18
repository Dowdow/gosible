# gosible

[![Tests](https://github.com/Dowdow/gosible/actions/workflows/test.yml/badge.svg)](https://github.com/Dowdow/gosible/actions/workflows/test.yml)

A reinterpretation of Ansible in Go for managing homeservers with a simple yaml configuration.

Under development. Only tested on Linux.

![Gif example](https://github.com/Dowdow/gosible/blob/main/demo/demo.gif?raw=true)

Regenerate it with `./demo/record.sh` (requires [vhs](https://github.com/charmbracelet/vhs)) — it runs against the fixtures in `demo/`, not a real machine.

## Installation

Download the binary from the [latest release](https://github.com/Dowdow/gosible/releases).

Move it to `/usr/local/bin/` or any of your `$PATH` directories.

### Build it yourself

```bash
git clone https://github.com/Dowdow/gosible.git
cd gosible
go mod download
make # Defaults to Build
make build # Create a build/gosible binary
make install # Move it to /usr/local/bin/gosible (works for update too)
make uninstall # Remove /usr/local/bin/gosible
make clean # Remove build/gosible
```

### Development

```bash
git clone https://github.com/Dowdow/gosible.git
cd gosible
go mod download
go run ./cmd/main.go config.yaml
```

## Usage

```bash
gosible config.yaml
```

## Config

The config is split in 3 categories :
- `inventory` the list and configuration of your machines.
- `actions` the repetitive actions.
- `tasks` a set of actions to execute in sequence.

If a `.env` file exists, it will be parsed and variables can be used in strings with the `env(ENV_VAR)` syntax, where `ENV_VAR` is the name of your environment variable.

```yaml
inventory:
  - id: machine1 # unique id for the machine
    name: Machine 1 # explicit name
    address: 192.168.1.10:22 # address:port
    users:
      - user: docker # the actual user
        ssh: /home/user/.ssh/id_rsa # the private ssh key...
        password: pass123 # ... or the user password
        become: pass456 # sudo password for privileged usage
        # or use environment variables:
        # password: env(PASSWORD)
        # become: env(SUDO_PASSWORD)
      # More users here
  # More machines here

actions:
  - id: unique.action.id # a unique id to identify the action in tasks
    name: Explicit action name # explicit action name
    type: shell # action type (see modules)
    args: {} # specific module args (see modules)
  # More actions here

tasks:
  - name: An automated task # explicit task name
    machines: # restrict the usage of the task to a specific machine and/or machine + user combo
      - machine1 # either a specific machine (use the id of the machine)...
      - machine1.docker # ...or a machine_id.user combo
    actions:
      - id: unique.action.id # either an action id from actions...
      - name: Explicit action name # ... or a specific action (same schema as actions)
        type: shell
        args: {}
      # More actions here
  # More tasks here
```
### Modules

List of modules to use in actions.

#### copy

Copy a file or a directory recursively to the machine. If `src` is relative, it will be from the config file path.

```yaml
name: A copy action
type: copy
args:
  src: relative/path/to/file.yml
  dest: /path/to/file.yml
  # or
  # src: /home/user/absolute/path/to/directory
  # dest: /path/to/directory
```

#### dir

Create dirs recursively. `mod` set the permissions (optional).

```yaml
name: A dir action
type: dir
args:
  path: /path/to/create
  mod: "644"
```

#### docker

Docker builds a `docker` image, saves it as a `.tar` file, and uploads it to the machine.

```yaml
name: A docker action
type: docker
args:
  src: ./image # directory of the Dockerfile
  dest: /path/to/the/image # must be a directory
  image: image:latest # docker image name
  tar: image-latest.tar # tar filename
  pull: true # always attempt to pull a newer version of the image
  clean: true # remove the local .tar file
```

#### file

Generate a file on the machine. `content` is an array and each element is a line in the file.

In this example, it generates a `.env` file, and the `env(ENV_VAR)` syntax can be used.

```yaml
name: A file action
type: file
args:
  dest: /home/user/app/.env
  content:
    - "TZ=Europe/Paris"
    - "CLIENT_ID=env(SERVICE_CLIENT_ID)"
    - "CLIENT_SECRET=env(SERVICE_CLIENT_SECRET)"
```

#### shell

Execute a shell command on the machine.

```yaml
name: A shell action
type: shell
args: "rm -rf /" # please don't
```
