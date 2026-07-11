Simple CLI tool for counting words in text files and directories that contains text files.

### Debug

Build and run docker container

```bash
./run -b
```

VS Code configuration is in `.vscode/launch.json`

```json
{
  "name": "Launch CLI Tool",
  "type": "go",
  "debugAdapter": "dlv-dap",
  "request": "attach",
  "mode": "remote",
  "substitutePath": [
    {
      "from": "${workspaceFolder}",
      "to": "/app"
    }
  ],
  "host": "127.0.0.1",
  "port": 40000,
  "showLog": true,
  "trace": "verbose"
}
```

Run Launch CLI Tool debug configuration.

### Run

Inside a docker container run

```bash
./w file_path/file_name1 file_path/file_name2 ...

or

./w *

or

./w some_directory/*
```

To build app without debug information run:

```bash
make build
```
