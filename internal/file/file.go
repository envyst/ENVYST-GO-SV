package file

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ENVYST-GO-SV/cryptography"
	"ENVYST-GO-SV/utilities"
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
		decryptedName, err := cryptography.FileNameDecrypt(file.Name(), password)
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
		decryptedName, err := cryptography.FileNameDecrypt(file.Name(), password)
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
func DeleteData() {
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
		utilities.ReturnPrompt()
		return
	}

	err := ListAndDelete(directory)
	if err != nil {
		fmt.Println("Error:", err)
	}
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
		decryptedName, err := cryptography.FileNameDecrypt(file.Name(), password)
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

	fmt.Println("Choose data type to add:")
	fmt.Println("1. Account\n2. Seed\n3. Private Key\n4. Other")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := scanner.Text()

	var directory, namePrompt, dataPrompt string
	switch choice {
	case "1":
		directory, namePrompt, dataPrompt = "accounts", "Enter Account Name", "Enter Username and Password"
	case "2":
		directory, namePrompt, dataPrompt = "seeds", "Enter Wallet Name", "Enter Seed"
	case "3":
		directory, namePrompt, dataPrompt = "private_keys", "Enter Wallet Name", "Enter Private Key"
	case "4":
		directory, namePrompt, dataPrompt = "others", "Enter Type Name", "Enter Data"
	default:
		return errors.New("invalid choice")
	}

	fmt.Println(namePrompt + ":")
	scanner.Scan()
	name := scanner.Text()

	exists, err := AccountNameExists(directory, name, password)
	if err != nil || exists {
		return errors.New("name already exists")
	}

	fmt.Println(dataPrompt + ":")
	scanner.Scan()
	data := scanner.Text()

	encryptedContent, err := cryptography.EncryptData(data, password)
	if err != nil {
		return err
	}

	fileName, err := cryptography.FileNameEncrypt(name, password)
	if err != nil {
		return err
	}

	return SaveToFile(directory, fileName, encryptedContent)
}
