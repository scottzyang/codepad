package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/chroma/v2/quick"
)

func main() {

	// get snippet

	err := quick.Highlight(os.Stdout, "func main() {}", "go", "terminal256", "monokai")
	if err != nil {
		fmt.Println("Error has occurred: ", err)
	}
}
