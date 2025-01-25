package utilities

import (
    "bufio"
    "fmt"
    "os"
)

func ClearScreen() {
    fmt.Print("\033[H\033[2J")
}

func ReturnPrompt() {
    fmt.Println("Press Enter to continue...")
    bufio.NewReader(os.Stdin).ReadBytes('\n')
}
