package file

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"ENVYST-GO-SV/internal/cryptography"
)

// ValidateSeed checks if the seed has 12 or 24 words.
func ValidateSeed(seed string) bool {
	words := strings.Fields(seed)
	return len(words) == 12 || len(words) == 24
}

// ValidatePrivateKey checks if the private key starts with "0x" and has 66 characters.
func ValidatePrivateKey(privateKey string) bool {
	return strings.HasPrefix(privateKey, "0x") && len(privateKey) == 66
}

// SaveToFile saves the given content to a file within the specified directory.
func SaveToFile(directory, fileName, content string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, os.ModePerm); err != nil {
			return err
		}
	}

	filePath := filepath.Join(directory, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// AccountNameExists checks if an account name already exists in the directory.
func AccountNameExists(directory, accountName, password string) (bool, error) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return false, nil
	}

	files, err := os.ReadDir(directory)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		decryptedName, err := cryptography.DecryptData(file.Name(), password)
		if err == nil && decryptedName == accountName {
			return true, nil
		}
	}

	return false, nil
}

// ListAndChoose lists entries in a directory and allows the user to select one.
func ListAndChoose(directory, password string) (string, error) {
	if password == "" {
		return "", errors.New("password not set")
	}

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return "", errors.New("no data available")
	}

	files, err := os.ReadDir(directory)
	if err != nil || len(files) == 0 {
		return "", errors.New("no data available")
	}

	fmt.Println("Available entries:")
	decryptedFiles := []string{}
	for i, file := range files {
		decryptedName, err := cryptography.DecryptData(file.Name(), password)
		if err == nil {
			decryptedFiles = append(decryptedFiles, file.Name())
			fmt.Printf("%d. %s\n", i+1, decryptedName)
		}
	}

	fmt.Print("Choose an entry by number: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := scanner.Text()
	idx, err := strconv.Atoi(choice)
	if err != nil || idx < 1 || idx > len(decryptedFiles) {
		return "", errors.New("invalid choice")
	}

	selectedFile := decryptedFiles[idx-1]
	fileContent, err := os.ReadFile(filepath.Join(directory, selectedFile))
	if err != nil {
		return "", err
	}

	decryptedContent, err := cryptography.DecryptData(string(fileContent), password)
	if err != nil {
		return "", err
	}

	fmt.Println("Details:")
	fmt.Println("------------------------------------")
	fmt.Println(decryptedContent)
	fmt.Println("------------------------------------")

	return decryptedContent, nil
}

// DeleteData prompts the user to choose a data type and deletes the selected item.
func DeleteData(password string) error {
	if password == "" {
		return errors.New("password not set")
	}

	options := `
1. Account
2. Wallet (Seed)
3. Wallet (Private Key)
4. Other
`
	fmt.Println(options)
	fmt.Print("Choose Account Type: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := scanner.Text()

	var directory string
	switch choice {
	case "1":
		directory = "accounts"
	case "2":
		directory = "seeds"
	case "3":
		directory = "private_keys"
	case "4":
		directory = "others"
	default:
		fmt.Println("Invalid option.")
		return nil // Return nil to avoid breaking the flow; option is invalid.
	}

	return ListAndDelete(directory, password)
}

// ListAndDelete lists entries and deletes the selected one.
func ListAndDelete(directory, password string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return errors.New("no data available")
	}

	files, err := os.ReadDir(directory)
	if err != nil || len(files) == 0 {
		return errors.New("no data available")
	}

	fmt.Println("Available entries:")
	decryptedFiles := []string{}
	for i, file := range files {
		decryptedName, err := cryptography.DecryptData(file.Name(), password)
		if err == nil {
			decryptedFiles = append(decryptedFiles, file.Name())
			fmt.Printf("%d. %s\n", i+1, decryptedName)
		}
	}

	fmt.Print("Choose to delete by number: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := scanner.Text()
	idx, err := strconv.Atoi(choice)
	if err != nil || idx < 1 || idx > len(decryptedFiles) {
		return errors.New("invalid choice")
	}

	fmt.Print("Are you sure? (y/n): ")
	scanner.Scan()
	if strings.ToLower(scanner.Text()) != "y" {
		return nil
	}

	selectedFile := decryptedFiles[idx-1]
	return os.Remove(filepath.Join(directory, selectedFile))
}

// AddData handles adding various types of data (Account, Seed, Private Key, etc.).
func AddData(password string) error {
	if password == "" {
		return errors.New("password not set")
	}

	// Display options and get choice
	fmt.Println("1. Account\n2. Seed\n3. Private Key\n4. Other")
	fmt.Print("Choose data type to add: ")
	choice, err := getUserInput()
	if err != nil {
		return err
	}

	// Prepare data based on choice
	directory, namePrompt, dataPrompt, processData := getDataDetails(choice)
	if directory == "" {
		return errors.New("invalid choice")
	}

	// Get name input
	fmt.Print(namePrompt + " : ")
	name, err := getUserInput()
	if err != nil {
		return err
	}

	// Check if name exists
	exists, err := AccountNameExists(directory, name, password)
	if err != nil || exists {
		return errors.New("name already exists")
	}

	// Collect content based on choice
	content, err := processData(choice, dataPrompt)
	if err != nil {
		return err
	}

	// Encrypt the content and name, then save to file
	encryptedContent, err := cryptography.EncryptData(content, password)
	if err != nil {
		return err
	}

	fileName, err := cryptography.EncryptData(name, password)
	if err != nil {
		return err
	}

	return SaveToFile(directory, fileName, encryptedContent)
}

// getUserInput handles user input and trims extra spaces.
func getUserInput() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}
	return "", errors.New("error reading input")
}

// getDataDetails returns details based on the user's choice.
func getDataDetails(choice string) (directory, namePrompt, dataPrompt string, processData func(string, string) (string, error)) {
	switch choice {
	case "1":
		return "accounts", "Enter Account Name", "Enter Data", processAccountData
	case "2":
		return "seeds", "Enter Wallet Name", "Enter Seed", processSimpleData
	case "3":
		return "private_keys", "Enter Wallet Name", "Enter Private Key", processSimpleData
	case "4":
		return "others", "Enter Type Name", "Enter Data", processOtherData
	default:
		return "", "", "", nil
	}
}

// processAccountData handles the data collection for account type.
func processAccountData(choice, dataPrompt string) (string, error) {
	var content string
	fmt.Println(dataPrompt)

	// Collect Username
	fmt.Print("Username : ")
	username, err := getUserInput()
	if err != nil {
		return "", err
	}
	content += "Username: " + username + "\n"

	// Collect Password
	fmt.Print("Password : ")
	password, err := getUserInput()
	if err != nil {
		return "", err
	}
	content += "Password: " + password + "\n"

	// Collect other data
	fmt.Print("Other Data? (y/n) : ")
	other, err := getUserInput()
	if err != nil {
		return "", err
	}

	for other == "y" {
		fmt.Print("Enter Data Key (Leave Blank to SKIP) : ")
		key, err := getUserInput()
		if err != nil {
			return "", err
		}
		if key != "" {
			content += key + ": "
			fmt.Print(key + " : ")
			value, err := getUserInput()
			if err != nil {
				return "", err
			}
			content += value + "\n"
		} else {
			other = "n"
		}
	}
	return content, nil
}

// processSimpleData handles simple data collection (seed, private key).
func processSimpleData(choice, dataPrompt string) (string, error) {
	fmt.Print(dataPrompt + " : ")
	data, err := getUserInput()
	if err != nil {
		return "", err
	}
	return data, nil
}

// processOtherData handles the data collection for "Other" type.
func processOtherData(choice, dataPrompt string) (string, error) {
	var content string
	for {
		fmt.Print("Enter Data Key (Leave Blank to SKIP) : ")
		key, err := getUserInput()
		if err != nil {
			return "", err
		}
		if key != "" {
			content += key + ": "
			fmt.Print(key + " : ")
			value, err := getUserInput()
			if err != nil {
				return "", err
			}
			content += value + "\n"
		} else {
			break
		}
	}
	return content, nil
}
