package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/fatih/color"
)

type CrudOption string

const (
	READ   CrudOption = "Read"
	WRITE  CrudOption = "Write"
	DELETE CrudOption = "Delete"
)

var crudOptions = []CrudOption{READ, WRITE, DELETE}

// option struct
type LanguageOption struct {
	Number       int
	LanguageName string
}

// snippet struct used for created file and snippet
type Snippet struct {
	Name           string
	Language       string
	SnippetContent []byte
}

func main() {
	var buffer bytes.Buffer

	// init codepad dir
	createNewCodepadDirectory()
	codePadDir := getHomeDir()

	reader := bufio.NewReader(os.Stdin)
	// get user CRUD selection
	switch selectedCrud := getUserCrudSelection(reader); selectedCrud {
	case READ:
		// get user selected languages
		selectedLanguage := getUserLanguage(reader, selectedCrud)
		langPath := codePadDir + "/" + selectedLanguage

		// get snippet selection
		snippet := getSnippetSelection(langPath, READ, reader)
		snippetPath := langPath + "/" + snippet

		// display snippet
		findAndDisplaySnippet(snippetPath, selectedLanguage)

	case WRITE:
		// get user selected language
		selectedLanguage := getUserLanguage(reader, selectedCrud)

		// get user code snippet
		userSnippet, snippetText := getUserSnippet(reader, buffer)

		newSnippet := Snippet{
			Name:           userSnippet,
			Language:       selectedLanguage,
			SnippetContent: snippetText,
		}

		// create new snippet within language directory
		createNewSnippet(newSnippet)
	case DELETE:
		// get user selected languages
		selectedLanguage := getUserLanguage(reader, selectedCrud)
		langPath := codePadDir + "/" + selectedLanguage

		// get snippet selection
		snippet := getSnippetSelection(langPath, DELETE, reader)
		snippetPath := langPath + "/" + snippet

		// delete snippet
		findAndDeleteSnippet(snippetPath)
	}
}

func getUserCrudSelection(reader *bufio.Reader) CrudOption {
	// display options
	for i, option := range crudOptions {
		formattedOption := fmt.Sprintf("%d. %s", i+1, option)
		fmt.Println(formattedOption)
	}

	for {
		blue := color.New(color.FgBlue)
		blue.Print("What would you like to do? (Input number): ")
		input, err := reader.ReadString('\n')

		input = strings.TrimSpace(input)
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		// verify input is number
		selectedOption, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if selectedOption == 0 || selectedOption > len(crudOptions) {
			fmt.Println("Option does not exist")
			continue
		}
		return crudOptions[selectedOption-1]
	}
}

func getUserLanguage(reader *bufio.Reader, selectedCrud CrudOption) string {
	// prompt user for language name
	blue := color.New(color.FgBlue)
	blue.Print("Select language: \n")

	// search all directories and return them as options.
	languages := getLanguageDirectories()

	// create options list
	optionsList := getOptionList(languages, selectedCrud)

	// display options list
	displayOptionsList(optionsList)

	// get user input
	userInput := getUserLanguageSelection(optionsList, reader)

	updatedSelectedLanguage := newLanguageInput(userInput.LanguageName, reader)

	userInput.LanguageName = updatedSelectedLanguage

	// wait for user input
	return userInput.LanguageName
}

func getUserLanguageSelection(optionsList []LanguageOption, reader *bufio.Reader) *LanguageOption {
	for {
		blue := color.New(color.FgBlue)
		blue.Print("Enter the selected language number: ")
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

		return inputOption
	}
}

func verifyUserInput(optionsList []LanguageOption, input string) (*LanguageOption, error) {
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

	return nil, errors.New("option does not exist")
}

func getUserSnippet(reader *bufio.Reader, buffer bytes.Buffer) (string, []byte) {
	// prompt user to add to existing or create new
	fmt.Println("Enter a title for the snippet (maximum 50 characters):")
	snippetTitle, err := reader.ReadString('\n')
	if err != nil {
		// Handle error
		fmt.Println("Error reading code snippet title", err)
	}

	// Trim any leading/trailing whitespace
	snippetTitle = strings.TrimSpace(snippetTitle)

	// Limit the input to a certain length, e.g., 50 characters
	maxChars := 50
	if len(snippetTitle) > maxChars {
		snippetTitle = snippetTitle[:maxChars]
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Paste code snippet (type 'done' on a new line and then hit enter to save):")

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

	return snippetTitle, buffer.Bytes()
}

func displayOptionsList(optionsList []LanguageOption) {
	for _, option := range optionsList {
		formattedOption := fmt.Sprintf("%d. %s", option.Number, option.LanguageName)
		fmt.Println(formattedOption)
	}
}

func getOptionList(languages []string, selectedCrud CrudOption) []LanguageOption {
	counter := 1
	var optionsList []LanguageOption

	for _, language := range languages {
		option := LanguageOption{
			Number:       counter,
			LanguageName: language,
		}

		optionsList = append(optionsList, option)
		counter += 1
	}

	if selectedCrud == WRITE {
		addNewLanguageOption := LanguageOption{
			Number:       counter,
			LanguageName: "Add a new language",
		}
		optionsList = append(optionsList, addNewLanguageOption)
	}
	return optionsList
}

// create new code pad folder for snippets
func createNewCodepadDirectory() {
	// Get the directory path
	codePadDir := getHomeDir()

	// Check if the directory exists, else create it
	directoryCheckMessage()
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
		directoryCheckSuccessMessage()
	}
}

// create language directory WIP refactor this to be more reusable with other directory making functions
func createNewLanguageDirectory(language string) {
	codePadDir := getHomeDir()
	capitalizedLanguage := capitalizeFirstLetter(language)
	languageDir := codePadDir + "/" + capitalizedLanguage

	// Check if the directory exists, else create it
	if _, err := os.Stat(languageDir); os.IsNotExist(err) {
		// if it doesn't exist, create it
		err := os.MkdirAll(languageDir, 0755)
		if err != nil {
			fmt.Println("Error creating language directory:", err)
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

func newLanguageInput(option string, reader *bufio.Reader) string {
	var updatedCapitalized string
	if option == "Add a new language" {
		fmt.Println("What language would you like to add?")
		input, err := reader.ReadString('\n')
		input = strings.TrimSpace(input) // Trim leading and trailing whitespace
		if err != nil {
			fmt.Println("Error reading input:", err)
			return "" // Return empty string on error
		}
		updatedCapitalized = capitalizeFirstLetter(input)
		createNewLanguageDirectory(updatedCapitalized)
		return updatedCapitalized
	}
	return option // Return empty string if the option is not "Add a new language"
}

func capitalizeFirstLetter(input string) string {
	if len(input) == 0 {
		return input
	}
	capitalized := strings.ToUpper(input[:1]) + input[1:]
	return capitalized
}

func createNewSnippet(newSnippet Snippet) {
	// get filepath
	filepath := getHomeDir()
	fmt.Println("Saved at" + filepath + "/" + newSnippet.Language + "/" + newSnippet.Name)

	// save the buffer to file
	file, err := os.Create(filepath + "/" + newSnippet.Language + "/" + newSnippet.Name)
	file.Write(newSnippet.SnippetContent)
	if err != nil {
		fmt.Println("Error creating file", err)
	}
}

// don't use this for write, read/delete only
func getSnippetSelection(path string, crudOption CrudOption, reader *bufio.Reader) string {
	var snippets []string

	dir, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening directory:", err)
	}
	defer dir.Close()

	// read all file names
	entries, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("Error reading directory contents:", err)
	}

	// Iterate over the entries
	for _, entry := range entries {
		// Check if the entry is a file
		if !entry.IsDir() {
			// append the file name
			snippets = append(snippets, entry.Name())
		}
	}

	// prompt for snippet selection
	if crudOption == READ {
		blue := color.New(color.FgBlue)
		blue.Println("Choose a snippet to read:")

	} else if crudOption == DELETE {
		blue := color.New(color.FgBlue)
		blue.Println("Choose a snippet to delete:")
	}

	for i, option := range snippets {
		formattedOption := fmt.Sprintf("%d. %s", i+1, option)
		fmt.Println(formattedOption)
	}

	for {
		blue := color.New(color.FgBlue)
		blue.Print("Input snippet selection number: ")
		input, err := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if err != nil {
			fmt.Println("Error reading input:", err)
			fmt.Println("Try again:")
			continue
		}
		// verify input is number
		selectedOption, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Not a number, try again:")
			continue
		}
		if selectedOption == 0 || selectedOption > len(crudOptions) {
			fmt.Println("Option does not exist, try again:")
			continue
		}
		return snippets[selectedOption-1]
	}
}

func findAndDisplaySnippet(path string, language string) {
	b, err := os.ReadFile(path) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	str := string(b) // convert content to a 'string'

	fmt.Println("")
	err = quick.Highlight(os.Stdout, str, language, "terminal256", "monokai")
	if err != nil {
		fmt.Println("Error has occurred: ", err)
	}
}

func findAndDeleteSnippet(path string) {
	e := os.Remove(path)
	if e != nil {
		fmt.Println("Failed to delete the file")
	}
	fmt.Println("Successfully deleted the snippet")
}

func directoryCheckMessage() {
	yellow := color.New(color.FgYellow)
	yellow.Printf("Checking if codepad directory exists...\n")
}

func directoryCheckSuccessMessage() {
	green := color.New(color.FgGreen)
	green.Printf("Directory already exists\n")
}
