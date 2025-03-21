package commands

import (
	"errors"
	"backend/stores"
	"fmt"
)

// no fields are needed for this command
type MOUNTED struct {}

func ParseMounted(tokens []string) (string, error) {
	fmt.Println("MOUNTED")
	// no arguments are needed
	if len(tokens) != 0 {
		return "", errors.New("no se esperaban argumentos")
	}

	stores.PrintMountedPartitions()

	return fmt.Sprintf("Se han mostrado las particiones montadas"), nil
	
}

