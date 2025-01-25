package logo

import (
	"fmt"
	"ENVYST-GO-SV/internal/utilities" // Import utilities package
)

// DrawLogo prints the application logo.
func DrawLogo() {
	logo := `
=====================================
   E  N  V  Y  S  T   (Secure Vault)
=====================================
`
	fmt.Println(logo)
}

// ShowMenu clears the screen, draws the logo, and displays the menu options.
func ShowMenu() {
	utilities.ClearScreen() // Use ClearScreen from utilities
	DrawLogo()

	menu := `
Select an option:
1. Reset Password
2. List Account
3. List Wallet (Seed)
4. List Wallet (Private Key)
5. List Other Data
6. Add
7. Delete
8. Google Drive Setup (Not Implemented)
9. Sync Accounts (Not Implemented)
------------------------------
`
	fmt.Println(menu)
}
