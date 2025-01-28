package main

import (
	"fmt"
	"os"
	"strings"
	"time"

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
			handeListData()
			utilities.ReturnPrompt()
		case "3":
			// Add data
			utilities.ClearScreen()
			handleAddData()
			utilities.ReturnPrompt()
		case "4":
			// Delete data
			utilities.ClearScreen()
			handleDeleteData()
			utilities.ReturnPrompt()
		case "5":
			// Example function for menu 8 (skipped functionality)
			utilities.ClearScreen()
			fmt.Println("Feature not implemented yet.")
			utilities.ReturnPrompt()
		case "6":
			// Example function for menu 9 (skipped functionality)
			utilities.ClearScreen()
			fmt.Println("Feature not implemented yet.")
			utilities.ReturnPrompt()
		case "exit", "q", "Q", "Exit":
			// Exit the application with 3-second delay
			fmt.Println("Exiting... Goodbye!")
			time.Sleep(1 * time.Second) // Delay for 3 seconds
			utilities.ClearScreen()
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
	fmt.Fscanln(os.Stdin, &pass)
	fmt.Println("Password set successfully.")
	return pass
}

// handeListData calls ListAndChoose from the file module.
func handeListData() {
	err := file.ListData(password)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
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
	} else {
		fmt.Println("Data deleted successfully!")
	}
}
