package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// convert a size and unit to bytes
func ConvertToBytes(size int, unit string) (int, error) {
	switch unit {
	case "K":
		return size * 1024, nil
	case "M":
		return size * 1024 * 1024, nil
	case "G":
		return size * 1024 * 1024 * 1024, nil
	case "T":
		return size * 1024 * 1024 * 1024 * 1024, nil
	default:
		return 0, errors.New("invalid unit")
	}
}

// alphabet array
var alphabet =	[]string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

// map to store the assignment of letters to the different paths
var pathToLetter = make(map[string]string)

var nextLetterIndex = 0

// get the letter assigned to a path
func GetLetter(path string) (string, error) {
	// Assign a letter to the path if it does not have one
	if _, exists := pathToLetter[path]; !exists {
		if nextLetterIndex < len(alphabet) {
			pathToLetter[path] = alphabet[nextLetterIndex]
			nextLetterIndex++
		} else {
			fmt.Println("Error: there are no more letters available to assign")
			return "", errors.New("no more letters available")
		}
	}

	return pathToLetter[path], nil
}

// createParentDirs crea las carpetas padre si no existen
func CreateParentDirs(path string) error {
	dir := filepath.Dir(path)
	// os.MkdirAll no sobrescribe las carpetas existentes, solo crea las que no existen
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error al crear las carpetas padre: %v", err)
	}
	return nil
}

// getFileNames obtiene el nombre del archivo .dot y el nombre de la imagen de salida
func GetFileNames(path string) (string, string) {
	dir := filepath.Dir(path)
	baseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	dotFileName := filepath.Join(dir, baseName+".dot")
	outputImage := path
	return dotFileName, outputImage
}