package main

import (
	"fmt"
	"strings"

	"ENVYST-GO-SV/internal/logo"
//	"ENVYST-GO-SV/internal/cryptography"
//	"ENVYST-GO-SV/internal/file"
	"ENVYST-GO-SV/internal/utilities"
//	"ENVYST-GO-SV/internal/gdrive"
)

func main() {
	utilities.ClearScreen()
	setupPassword()

	for {
		logo.ShowMenu()
		fmt.Print("Enter your choice: ")
		var choice string
		fmt.Scanln(&choice)
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			utilities.ClearScreen()
			setupPassword()
                        utilities.ReturnPrompt()
		case "2":
			utilities.ClearScreen()
			listAndChoose("accounts")
                        utilities.ReturnPrompt()
		case "3":
			utilities.ClearScreen()
			listAndChoose("seeds")
                        utilities.ReturnPrompt()
		case "4":
			utilities.ClearScreen()
			listAndChoose("private_keys")
                        utilities.ReturnPrompt()
		case "5":
			utilities.ClearScreen()
			listAndChoose("others")
                        utilities.ReturnPrompt()
		case "6":
			utilities.ClearScreen()
			addData()
                        utilities.ReturnPrompt()
		case "7":
			utilities.ClearScreen()
			deleteData()
                        utilities.ReturnPrompt()
		case "8":
			utilities.ClearScreen()
			setupGoogleCredentials()
                        utilities.ReturnPrompt()
		case "9":
			utilities.ClearScreen()
			syncAccount()
                        utilities.ReturnPrompt()
		case "exit", "q", "Q", "Exit":
			fmt.Println("Exiting... Goodbye!")
			utilities.ClearScreen()
			return
		default:
			fmt.Println("Invalid option. Try again.")
                        utilities.ReturnPrompt()
		}
	}
}

// Example functions
func setupPassword() {
	fmt.Println("Setting up password...")
}

func listAndChoose(category string) {
	fmt.Printf("Listing and choosing from category: %s\n", category)
}

func addData() {
	fmt.Println("Adding data...")
}

func deleteData() {
	fmt.Println("Deleting data...")
}

func setupGoogleCredentials() {
	fmt.Println("Setting up Google credentials...")
}

func syncAccount() {
	fmt.Println("Syncing account...")
}
