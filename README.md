# enumgen

`enumgen` is a specialized tool for generating Go enum types from a simple, declarative language. It eliminates repetitive boilerplate code for enum types in Go, providing type safety and useful methods.

## Overview

Go doesn't have native enum support, making it tedious to implement consistent, type-safe enums with helpful methods. `enumgen` solves this by providing a simple language to define enums and automatically generates the corresponding Go code.

## Features

- Simple, declarative syntax for defining enums
- Supports multiple value types (string, int, etc.)
- Supports key-value pairs
- Generates type-safe enum implementations
- Generates helper methods (String(), IsValid(), etc.)
- Integrates with `go generate`

## Installation

```bash
go install github.com/kkumar-gcc/enumgen@latest
```

## Usage

### Basic Usage with go generate

Add a comment directive to your Go file:

```go
//go:generate enumgen -file=enums.edl -output=generated_enums.go
```

Then run:

```bash
go generate ./...
```

### EDL File Syntax

Create an `.edl` (Enum Definition Language) file with your enum definitions:

```
// Status represents operation result status
enum Status [int]:
    SUCCESS = 0,
    WARNING = 1,
    ERROR = 2;

// Color defines standard colors
enum Color [string]:
    RED = "red",
    GREEN = "green",
    BLUE = "blue";
```

### Enum Definition Types

#### Simple Enum (without type)

```
enum Direction:
    NORTH,
    EAST,
    SOUTH,
    WEST;
```

#### Enum with Single Type

```
enum Status [int]:
    SUCCESS = 0,
    WARNING = 1,
    ERROR = 2;
```

#### Enum with Key-Value Pairs

```
enum Day [string, string]:
    MONDAY = "Monday":"Mon",
    TUESDAY = "Tuesday":"Tue",
    WEDNESDAY = "Wednesday":"Wed";
```

### Command Line Options

```
Usage:
  enumgen [flags]

Flags:
  -file string   Input EDL file (default "example.edl")
  -output string Output file for generated Go code
  -ast           Generate AST visualization (requires Graphviz)
```

## Grammar

For the complete grammar definition, see [grammar.md](grammar.md).

## Examples

See the `example.edl` file for reference examples.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
```