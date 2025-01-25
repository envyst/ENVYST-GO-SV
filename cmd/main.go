package main

import (
	"fmt"
	"strings"
	"time"

	"ENVYST-GO-SV/internal/cryptography"
	"ENVYST-GO-SV/internal/file"
	"ENVYST-GO-SV/internal/logo"
	"ENVYST-GO-SV/internal/utilities"
)

var password string

func main() {
	utilities.ClearScreen()
	password = setupPassword()
	for {
		logo.ShowMenu()
		fmt.Print("Enter your choice: ")
		var choice string
		fmt.Scanln(&choice)
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			// Reset password
			utilities.ClearScreen()
			password = setupPassword()
			utilities.ReturnPrompt()
		case "2":
			// List and choose account
			utilities.ClearScreen()
			handleListAndChoose("accounts")
			utilities.ReturnPrompt()
		case "3":
			// List and choose wallet seed
			utilities.ClearScreen()
			handleListAndChoose("seeds")
			utilities.ReturnPrompt()
		case "4":
			// List and choose wallet private key
			utilities.ClearScreen()
			handleListAndChoose("private_keys")
			utilities.ReturnPrompt()
		case "5":
			// List and choose other data
			utilities.ClearScreen()
			handleListAndChoose("others")
			utilities.ReturnPrompt()
		case "6":
			// Add data
			utilities.ClearScreen()
			handleAddData()
			utilities.ReturnPrompt()
		case "7":
			// Delete data
			utilities.ClearScreen()
			handleDeleteData()
			utilities.ReturnPrompt()
		case "8":
			// Example function for menu 8 (skipped functionality)
			utilities.ClearScreen()
			fmt.Println("Feature not implemented yet.")
			utilities.ReturnPrompt()
		case "9":
			// Example function for menu 9 (skipped functionality)
			utilities.ClearScreen()
			fmt.Println("Feature not implemented yet.")
			utilities.ReturnPrompt()
		case "exit", "q", "Q", "Exit":
			// Exit the application with 3-second delay
			utilities.ClearScreen()
			fmt.Println("Exiting... Goodbye!")
			time.Sleep(3 * time.Second) // Delay for 3 seconds
			return
		default:
			// Invalid option
			fmt.Println("Invalid option. Try again.")
			utilities.ReturnPrompt()
		}
	}
}

// setupPassword initializes or updates the password.
func setupPassword() string {
	fmt.Print("Enter your password: ")
	var pass string
	fmt.Fscanln(os.Stdin, &password)
	fmt.Println("Password set successfully.")
	return pass
}

// handleListAndChoose calls ListAndChoose from the file module.
func handleListAndChoose(directory string) {
	data, err := file.ListAndChoose(directory, password)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Selected data:", data)
}

// handleAddData calls AddData from the file module.
func handleAddData() {
	err := file.AddData(password)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Data added successfully!")
}

// handleDeleteData calls DeleteData from the file module.
func handleDeleteData() {
	err := file.DeleteData(password)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Data deleted successfully!")
}
