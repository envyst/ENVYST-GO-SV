package logo

import (
	"ENVYST-GO-SV/internal/utilities" // Import utilities package
	"fmt"
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
2. List Accounts
3. Add
4. Delete
5. Google Drive Setup (Not Implemented)
6. Sync Accounts (Not Implemented)
------------------------------
`
	fmt.Println(menu)
}
