# goat (Go Application Template)

A command-line interface (CLI) tool with an interactive UI to quickly bootstrap Go application templates.

## Description

`goat` is a modern CLI tool designed to help developers quickly generate and bootstrap Go application templates. It features an interactive terminal user interface built with Bubble Tea, making project creation intuitive and user-friendly.

## Features

- Interactive terminal UI for project configuration
- Creates Fiber framework projects with proper structure
- Step-by-step project setup wizard
- Generates complete project boilerplate (WIP)

## Installation

### Prerequisites

- Go 1.23.4 or later
- Make (optional, for using Makefile commands)

### Building from Source

1. Clone the repository:

```bash
git clone https://github.com/smilepakawat/goat.git
cd goat
```

2. Build the project:
```bash
go build -o build/goat
```

## Usage

### Creating a Fiber Project

Run the following command and follow the interactive prompts:

```bash
# Create a project with fiber framework
goat create-fiber

# Create a project with gin gonic framework
goat create-gin
```

You will be asked to provide:

1. Project name
2. Module name (Go module path)

After completion, your Fiber project will be created with all necessary boilerplate code.

### Next Steps

Once your project is created:

```bash
cd your-project-name
go run main.go
```

## Development

### Dependencies

- github.com/spf13/cobra - CLI framework
- github.com/charmbracelet/bubbletea - Terminal UI framework
- github.com/charmbracelet/bubbles - UI components

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the terms of the LICENSE file included in the repository.
