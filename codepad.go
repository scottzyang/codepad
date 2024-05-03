package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/alecthomas/chroma/v2/quick"
)

func main() {
	var buffer bytes.Buffer

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "done" {
			break
		}
		buffer.WriteString(line + "\n")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	err := quick.Highlight(os.Stdout, buffer.String(), "go", "terminal256", "monokai")
	if err != nil {
		fmt.Println("Error has occurred: ", err)
	}
}
