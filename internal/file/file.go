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

const parentDir = "data/"

// ValidateSeed checks if the seed has 12 or 24 words.
func ValidateSeed(seed string) bool {
	words := strings.Fields(seed)
	return len(words) == 12 || len(words) == 24
}

// ValidatePrivateKey checks if the private key starts with "0x" and has 66 characters.
func ValidatePrivateKey(privateKey string) bool {
	return strings.HasPrefix(privateKey, "0x") && len(privateKey) == 66
}

// EnsureDirectory ensures that a directory exists, creating it if necessary.
func EnsureDirectory(directory string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return os.MkdirAll(directory, os.ModePerm)
	}
	return nil
}

// SaveToFile saves the given content to a file within the specified directory.
func SaveToFile(directory, fileName, content string) error {
	if err := EnsureDirectory(directory); err != nil {
		return fmt.Errorf("failed to ensure directory %s: %w", directory, err)
	}

	filePath := filepath.Join(directory, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// DeleteData allows deleting a specific type of data (Account, Seed, etc.) after listing it.
func DeleteData(password string) error {
	return PerformDataOperation(password, func(directory, password string) error {
		return ListAndPerform(directory, password, func(filePath, fileName string) error {
			fmt.Print("Are you sure? (y/n): ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			if strings.ToLower(scanner.Text()) != "y" {
				return nil
			}
			fmt.Print(fileName + " ")
			return os.Remove(filePath)
		})
	})
}

// ListData allows listing data of a specific type (Account, Seed, etc.) without performing actions.
func ListData(password string) error {
	return PerformDataOperation(password, func(directory, password string) error {
		return ListAndPerform(directory, password, func(filePath, fileName string) error {
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			decryptedContent, err := cryptography.DecryptData(string(fileContent), password)
			if err != nil {
				return err
			}

			fmt.Println("Details:")
			fmt.Println("------------------------------------")
			fmt.Println(fileName)
			fmt.Println("------------------------------------")
			fmt.Println(decryptedContent)
			fmt.Println("------------------------------------")

			return nil
		})
	})
}

// ListAndPerform lists entries in a directory and allows the user to perform an action on the selected entry.
func ListAndPerform(directory, password string, action func(string, string) error) error {
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

	fmt.Print("Choose an entry by number: ")
	choice := getInput()
	idx, err := strconv.Atoi(choice)
	if err != nil || idx < 1 || idx > len(decryptedFiles) {
		return errors.New("invalid choice")
	}

	selectedFile := decryptedFiles[idx-1]
	decryptedFileName, err := cryptography.DecryptData(selectedFile, password)
	if err != nil {
		return err
	}
	return action(filepath.Join(directory, selectedFile), decryptedFileName)
}

// PerformDataOperation performs an operation (e.g., list or delete) on a data type.
func PerformDataOperation(password string, operation func(string, string) error) error {
	options := `
1. Account
2. Wallet (Seed)
3. Wallet (Private Key)
4. Other
`
	fmt.Println(options)
	fmt.Print("Choose Account Type: ")
	choice := getInput()

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
		return errors.New("invalid option")
	}

	return operation(filepath.Join(parentDir, directory), password)
}

// AddData handles adding various types of data (Account, Seed, Private Key, etc.).
func AddData(password string) error {
	dataOptions := map[string]struct {
		Dir         string
		NamePrompt  string
		DataPrompt  string
		ProcessFunc func() (string, error)
	}{
		"1": {"accounts", "Enter Account Name", "Enter Data", processAccountData},
		"2": {"seeds", "Enter Wallet Name", "Enter Seed", processSimpleData},
		"3": {"private_keys", "Enter Wallet Name", "Enter Private Key", processSimpleData},
		"4": {"others", "Enter Type Name", "Enter Data", processOtherData},
	}

	fmt.Println("1. Account\n2. Seed\n3. Private Key\n4. Other")
	fmt.Print("Choose data type to add: ")
	choice := getInput()

	option, valid := dataOptions[choice]
	if !valid {
		return errors.New("invalid choice")
	}

	fmt.Print(option.NamePrompt + ": ")
	accountName := getInput()

	exists, err := AccountNameExists(filepath.Join(parentDir, option.Dir), accountName, password)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("name already exists")
	}

	dataContent, err := option.ProcessFunc()
	if err != nil {
		return err
	}

	if dataContent != "" {
		encryptedContent, err := cryptography.EncryptData(dataContent, password)
		if err != nil {
			return err
		}

		encryptedFileName, err := cryptography.EncryptData(accountName, password)
		if err != nil {
			return err
		}

		return SaveToFile(filepath.Join(parentDir, option.Dir), encryptedFileName, encryptedContent)
	}

	return nil
}

// getInput handles user input
func getInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// processAccountData handles the data collection for account type.
func processAccountData() (string, error) {
	var content strings.Builder

	fmt.Print("Username: ")
	username := getInput()
	content.WriteString("Username: " + username + "\n")

	fmt.Print("Password: ")
	password := getInput()
	content.WriteString("Password: " + password + "\n")

	for {
		fmt.Print("Enter Data Key (Leave Blank to Skip): ")
		key := getInput()
		if key == "" {
			break
		}
		fmt.Print(key + ": ")
		value := getInput()
		content.WriteString(key + ": " + value + "\n")
	}

	return content.String(), nil
}

// processSimpleData handles simple data collection (seed, private key).
func processSimpleData() (string, error) {
	fmt.Print("Enter Data: ")
	data := getInput()
	if !ValidatePrivateKey(data) && !ValidateSeed(data) {
		return "", errors.New("invalid data format")
	}
	var content strings.Builder
	content.WriteString(data)
	return content.String(), nil
}

// processOtherData handles the data collection for "Other" type.
func processOtherData() (string, error) {
	var content strings.Builder

	for {
		fmt.Print("Enter Data Key (Leave Blank to Skip): ")
		key := getInput()
		if key == "" {
			break
		}
		fmt.Print(key + ": ")
		value := getInput()
		content.WriteString(key + ": " + value + "\n")
	}

	return content.String(), nil
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
