[![Go Report Card](https://goreportcard.com/badge/github.com/scottzyang/codepad)](https://goreportcard.com/report/github.com/scottzyang/codepad)

# CodePad - Command Line Code Snippet Manager

CodePad is a command-line tool written in Go for managing code snippets. It allows you to organize your code snippets by language and perform CRUD (Create, Read, Update, Delete) operations on them easily.

## Features

- **CRUD Operations:** Perform Create, Read, Update, and Delete operations on code snippets.
- **Syntax Highlighting:** Code snippets are displayed with syntax highlighting for better readability through the `.
- **Interactive Interface:** Easy-to-use command-line interface for managing code snippets.
- **Colorful Output:** Utilizes the `github.com/fatih/color` package for colorful console output.

## Tech Stack

- **Language:** Go
- **Packages:** `bufio`, `bytes`, `errors`, `fmt`, `os`, `strconv`, `strings`, `github.com/alecthomas/chroma/v2/quick`, `github.com/fatih/color`

## Usage
1. Build the program
    ```bash
    go build
    ```
2. Run the program
    ```bash
    ./codepad
    ```
3. Select an option
   - Choose one of the CRUD options presented.
   - Follow the prompts to select a language and perform actions on snippets.

## Example
```bash
./codepad
Checking if codepad directory exists...
Directory already exists
1. Read
2. Write
3. Delete
What would you like to do? (Input number): 1
Select language:
1. Go
2. Python
3. Ruby
Enter the selected language number: 2
Choose a snippet to read:
1. For loops
2. Merge two lists into a dictionary
Input snippet selection number: 1

fruits = ["apple", "banana", "cherry"]
for x in fruits:
  print(x)
```

## Contributing

Contributions are welcome! If you have any ideas for improvements or find any issues, feel free to open an issue or submit a pull request.

## Improvements

- Modular functions: maintain DRY code
  - Displaying options
  - Reading directories
- File structure (separate out by context)
- Have top level variables within main to reduce reinstantiation
  - example, `codePadDir := getHomeDir()` seen multiple times; can be called once
- Consolidate where errors are being handled

