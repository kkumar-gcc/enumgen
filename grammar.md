# Enumgen Grammar

Enumgen uses a custom grammar format to define enums. This document outlines the grammar rules and provides examples for reference.

## Grammar Rules

The grammar is defined in Extended Backus-Naur Form (EBNF) notation:

```ebnf
SourceFile       ::= { Definition } EOF ;

Definition       ::= EnumDefinition | Comment ;

EnumDefinition   ::= { Comment } 'enum' Identifier [ TypeSpec ] ':' MemberList ;

TypeSpec         ::= '[' Type { ',' Type } ']' ;

Type             ::= Identifier { '.' Identifier } ;

MemberList       ::= { MemberDefinition } ;

MemberDefinition ::= { Comment } Identifier [ MemberAssignment ] [ Terminator ] ;

MemberAssignment ::= '=' ( Literal | KeyValue ) ;

KeyValue         ::= Literal ':' Literal ;

Literal          ::= INT | FLOAT | CHAR | STRING | Identifier ;

Comment          ::= '//' .+ ;

Terminator       ::= ',' | ';' ;
```

## Grammar Explanation

### Enum Definition

An enum definition begins with the `enum` keyword, followed by the enum name (identifier) and optional type specifications. The enum members are listed after a colon `:`.

- **Identifier**: The name of the enum (e.g., `Color`, `Status`)
- **TypeSpec**: Optional type specifications in square brackets (e.g., `[string, int]`)
- **MemberList**: A list of enum members

### Type Specifications

Type specifications define the types associated with enum members. They are enclosed in square brackets and can include multiple types separated by commas.

For example: `[string, int]` specifies that the enum uses both string and int types.

### Enum Members

Enum members consist of an identifier (the member name) and an optional value assignment. The value assignment can be a simple literal or a key-value pair.

Members can be terminated by either a comma or a semicolon.

### Comments

Comments start with `//` and continue to the end of the line. They can be placed before an enum definition or before enum members.

## Examples

### Basic Enum with String Values

```
// Color enum defines basic colors
enum Color [string]:
    RED = "red",
    GREEN = "green",
    BLUE = "blue";
```

### Enum with Key-Value Pairs

```
// Days enum with both name and abbreviation
enum Day [string, string]:
    MONDAY = "Monday":"Mon",
    TUESDAY = "Tuesday":"Tue",
    WEDNESDAY = "Wednesday":"Wed",
    THURSDAY = "Thursday":"Thu",
    FRIDAY = "Friday":"Fri",
    SATURDAY = "Saturday":"Sat",
    SUNDAY = "Sunday":"Sun";
```

### Enum with Integer Values

```
// Status enum with integer codes
enum Status [int]:
    SUCCESS = 0,
    WARNING = 1,
    ERROR = 2,
    CRITICAL = 3;
```

### Enum without Explicit Values (Implicit Sequential Numbering)

```
// Direction enum without explicit values
enum Direction:
    NORTH,
    EAST,
    SOUTH,
    WEST;
```
