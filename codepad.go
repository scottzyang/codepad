package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	var buffer bytes.Buffer

	reader := bufio.NewReader(os.Stdin)

	// prompt user fo rlanguage name
	fmt.Println("Select language:")

	// search all directories and return them as options.
	counter := 1
	languages := get_language_directories()
	for _, language := range languages {
		option := fmt.Sprintf("%d. %s", counter, language)
		fmt.Println(option)
		counter += 1
	}

	// wait for user input
	selected_language, _ := reader.ReadString('\n')
	fmt.Println(selected_language)

	// use selected language to navigate to folder.

	// prompt user to add to existing or create new
	fmt.Println("Enter name of snippet:")
	snippet, _ := reader.ReadString('\n')

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Paste code snippet:")

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

	// save the buffer to file
	file, err := os.Create("./" + snippet)
	file.Write(buffer.Bytes())
	if err != nil {
		fmt.Println("Error creating file", err)
	}

	// // reading
	// // err := quick.Highlight(os.Stdout, buffer.String(), "go", "terminal256", "monokai")
	// // if err != nil {
	// // 	fmt.Println("Error has occurred: ", err)
	// // }
}

// create new code pad folder for snippets
func create_new_codepad_directory() {
	// Get the directory path
	codePadDir := get_home_dir()

	// Check if the directory exists, else create it
	if _, err := os.Stat(codePadDir); os.IsNotExist(err) {
		// if it doesn't exist, create it
		err := os.MkdirAll(codePadDir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		fmt.Println("Directory created successfully.")
	} else if err != nil {
		fmt.Println("Error checking directory existence:", err)
	} else {
		fmt.Println("Directory already exists.")
	}
}

// create language directory WIP refactor this to be more reusable with other directory making functions
func create_new_language_directory(language string) {
	codePadDir := get_home_dir()
	languageDir := codePadDir + "/" + language

	// Check if the directory exists, else create it
	if _, err := os.Stat(languageDir); os.IsNotExist(err) {
		// if it doesn't exist, create it
		err := os.MkdirAll(languageDir, 0755)
		if err != nil {
			fmt.Println("Error creating language directory:", err)
			return
		}
		fmt.Println("Language directory created successfully")
	} else if err != nil {
		fmt.Println("Error checking if language directory existence:", err)
	} else {
		fmt.Println("Language directory already exists.")
	}
}

// get the home dir for code pad
func get_home_dir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user's home directory:", err)
	}
	return homeDir + "/codepad"
}

// search for language directories within codepad
func get_language_directories() []string {
	var languageNames []string

	codePadDir := get_home_dir()

	dir, err := os.Open(codePadDir)
	if err != nil {
		fmt.Println("Error opening directory:", err)
	}
	defer dir.Close()

	// Read all entries in the directory
	entries, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("Error reading directory contents:", err)
	}

	// Iterate over the entries
	for _, entry := range entries {
		// Check if the entry is a directory
		if entry.IsDir() {
			// Print the directory name
			languageNames = append(languageNames, entry.Name())
		}
	}
	return languageNames
}
