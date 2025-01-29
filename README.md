# ENVYST Secure Vault (Local CLI)

A local encrypted vault for managing credentials and sensitive accounts. Data is stored securely on your machine with AES encryption.  
**Note:** Google Drive Sync and Cloud Sync are not yet implemented.


## Installation

### Pre-built Binary
1. Download the latest `envyst_amd64` binary from [Releases](https://github.com/envyst/ENVYST-GO-SV/releases).
2. Make it executable:
   ```bash
   chmod +x envyst_amd64 # Or chmod +x envyst_arm64
   ```

### Build from Source (Requires Go 1.20+)
```bash
git clone https://github.com/envyst/ENVYST-GO-SV
cd ENVYST-GO-SV
GOOS=linux GOARCH=amd64 go build -o bin/envyst_amd64 cmd/main.go
# Or
# GOOS=linux GOARCH=arm64 go build -o bin/envyst_arm64 cmd/main.go
```

---

## First Run & Setup
1. Start the app:
   ```bash
   ./bin/envyst_amd64
   ```
2. **Set a Master Password** when prompted (required for encryption/decryption).

---

## Menu Options

```
Select an option:
1. Reset Password
2. List Accounts
3. Add Account
4. Delete Account
5. Google Drive Setup (Not Implemented)
6. Sync Accounts (Not Implemented)
```

### 1. Reset Password
- Change your master password. Requires current password verification.

### 2. List Accounts
- View all stored accounts (ex) :
  ```
    ------------------------------------
    Username: test
    Password: test
    test: test

    ------------------------------------
  ```

### 3. Add Account
- Add new credentials:
  1. Enter **Service Name** (e.g., "Netflix")
  2. Enter **Username/Email**
  3. Enter **Password** 

### 4. Delete Account
- Remove an account by its **ID** (visible in "List Accounts").

---

## Data Storage
- Encrypted data stored at: `data/`
- **Master Password** is required to decrypt data. **Losing it = Permanent data loss!**

---

## Security
- 🔒 AES-256 encryption for stored data
- 🚫 No telemetry or external data transmission
- 📁 Local storage only (no cloud sync unless manually implemented)

---

## Example Workflow

```bash
# Start the vault
./bin/envyst_amd64

# First-time setup: Set master password
> Enter master password: ******

# Add a new account
> Select option 3
> Service: GitHub
> Username: user@email.com
> Password: ********

# List accounts
> Select option 2
------------------------------------
GitHub
Username: user@email.com
Password: ********

------------------------------------

# Delete account ID 1
> Select option 4
> Enter account ID to delete: 1
```

---

## Limitations
- ❌ Google Drive backup/sync not implemented
- ❌ Multi-device sync unavailable
- ❌ No password strength enforcement

---

**Warning:** This is a development-stage tool. Always keep backups of your `data/`!
```
