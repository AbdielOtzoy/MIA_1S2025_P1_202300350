package commands

import (
	"errors"
	"backend/stores"
	"fmt"
)

// no fields are needed for this command
type MOUNTED struct {}

func ParseMounted(tokens []string) (*MOUNTED, error) {
	fmt.Println("MOUNTED")
	// no arguments are needed
	if len(tokens) != 0 {
		return nil, errors.New("no se esperaban argumentos")
	}

	stores.PrintMountedPartitions()

	return &MOUNTED{}, nil
	
}

