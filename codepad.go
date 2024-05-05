package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// option struct
type Option struct {
	Number       int    // Option number
	LanguageName string // Option text
}

// snippet struct used for created file and snippet
type Snippet struct {
	Name           string
	Language       string
	SnippetContent string
}

func main() {
	var buffer bytes.Buffer

	reader := bufio.NewReader(os.Stdin)

	// get user selected language
	selectedLanguage := getUserLanguage(reader)
	fmt.Println(selectedLanguage)

	// get user code snippet
	userSnippet := getUserSnippet(reader)
	fmt.Println(userSnippet)

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
	file, err := os.Create("./" + userSnippet)
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

func getUserLanguage(reader *bufio.Reader) string {
	// prompt user for language name
	fmt.Println("Select language:")

	// search all directories and return them as options.
	languages := getLanguageDirectories()

	// create options list
	optionsList := getOptionList(languages)

	// display options list
	displayOptionsList(optionsList)

	// get user input
	userInput := getUserLanguageSelection(optionsList, reader)

	// wait for user input
	return userInput.LanguageName
}

func getUserLanguageSelection(optionsList []Option, reader *bufio.Reader) *Option {
	for {
		fmt.Print("Enter the selected language number: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		// verify user input exists as option
		inputOption, err := verifyUserInput(optionsList, input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// if input not verified reprompt
		return inputOption
	}
}

func verifyUserInput(optionsList []Option, input string) (*Option, error) {
	input = strings.TrimSpace(input)

	// verify input exists as option
	selectedOption, err := strconv.Atoi(input)
	if err != nil {
		return nil, err
	}
	for _, option := range optionsList {
		if option.Number == selectedOption {
			return &option, nil
		}
	}

	return nil, errors.New("Option does not exist")
}

func getUserSnippet(reader *bufio.Reader) string {
	// prompt user to add to existing or create new
	fmt.Println("Enter a title for the snippet (maximum 50 characters):")
	snippet, err := reader.ReadString('\n')
	if err != nil {
		// Handle error
		fmt.Println("Error reading code snippet title", err)
	}

	// Trim any leading/trailing whitespace
	snippet = strings.TrimSpace(snippet)

	// Limit the input to a certain length, e.g., 50 characters
	maxChars := 50
	if len(snippet) > maxChars {
		snippet = snippet[:maxChars]
	}

	return snippet
}

func displayOptionsList(optionsList []Option) {
	for _, option := range optionsList {
		formattedOption := fmt.Sprintf("%d. %s", option.Number, option.LanguageName)
		fmt.Println(formattedOption)
	}
}

func getOptionList(languages []string) []Option {
	counter := 1
	var optionsList []Option

	for _, language := range languages {
		option := Option{
			Number:       counter,
			LanguageName: language,
		}

		optionsList = append(optionsList, option)
		counter += 1
	}

	addNewLanguageOption := Option{
		Number:       counter,
		LanguageName: "Add a new language",
	}

	optionsList = append(optionsList, addNewLanguageOption)
	return optionsList
}

// create new code pad folder for snippets
func createNewCodepadDirectory() {
	// Get the directory path
	codePadDir := getHomeDir()

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
func createNewLanguageDirectory(language string) {
	codePadDir := getHomeDir()
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
func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user's home directory:", err)
	}
	return homeDir + "/codepad"
}

// search for language directories within codepad
func getLanguageDirectories() []string {
	var languageNames []string

	codePadDir := getHomeDir()

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
