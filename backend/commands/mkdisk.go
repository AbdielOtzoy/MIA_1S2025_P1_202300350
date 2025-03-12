package commands

import (
	structures "backend/structures"
	utils "backend/utils"
	"fmt"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"strconv"
	"time"
)

type MKDISK struct {
	Size int
	unit string
	fit string  // types: FF, BF, WF
	path string
}

func ParseMkdisk(tokens []string) (*MKDISK, error) {
	cmd := &MKDISK{} // create the mkdisk command

	args := strings.Join(tokens, " ") // join the tokens to get the arguments

	// regular expression to get the size and unit
	re := regexp.MustCompile(`-size=\d+|-unit=[kKmM]|-fit=[bBfFwW]{2}|-path="[^"]+"|-path=[^\s]+`)

	matches := re.FindAllString(args, -1)

	for _, match := range matches {
		// dibide the argument in the key and value
		kv := strings.SplitN(match, "=", 3)
		if len(kv) != 2 {
			return nil, errors.New("invalid argument")
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key {
			case "-size":
				size, err := strconv.Atoi(value) // convert the value to integer
				if err != nil {
					return nil, errors.New("invalid size")
				}
				cmd.Size = size
			case "-unit":
				value = strings.ToUpper(value)
				if value != "K" && value != "M" {
					return nil, errors.New("invalid unit")
				}
				cmd.unit = value
			case "-fit":
				value = strings.ToUpper(value)
				if value != "FF" && value != "BF" && value != "WF" {
					return nil, errors.New("invalid fit")
				}
				cmd.fit = value
			case "-path":
				if value == "" {
					return nil, errors.New("invalid path")
				}
				cmd.path = value
			default: 
				return nil, errors.New("invalid argument") // invalid argument

		}
	}

	if cmd.Size == 0 {
		return nil, errors.New("size is required")
	}

	if cmd.path == "" {
		return nil, errors.New("path is required")
	}

	if cmd.unit == "" {
		cmd.unit = "M"
	}

	if cmd.fit == "" {
		cmd.fit = "FF"
	}

	err := commandMkdisk(cmd)
	if err != nil {
		fmt.Println(err)
	}

	return cmd, nil
}

func commandMkdisk(mkdisk *MKDISK) error {
	sizeBytes, err := utils.ConvertToBytes(mkdisk.Size, mkdisk.unit) // convert the size to bytes
	if err != nil {
		fmt.Println("error converting size to bytes:", err)
		return err
	}

	// create the disk
	err = createDisk(mkdisk, sizeBytes)
	if err != nil {
		fmt.Println("error creating disk:", err)
		return err
	}

	err = createMBR(mkdisk, sizeBytes)
	if err != nil {
		fmt.Println("error creating MBR:", err)
		return err
	}

	return nil
}

func createDisk(mkdisk *MKDISK, sizeBytes int) error {
	err := os.MkdirAll(filepath.Dir(mkdisk.path), os.ModePerm) // create the directory
	if err != nil {
		fmt.Println("error creating directory:", err)
		return err
	}

	file, err := os.Create(mkdisk.path) // create the file
	if err != nil {
		fmt.Println("error creating file:", err)
		return err
	}
	defer file.Close()

	buffer := make([]byte, 1024*1024) // 1MB buffer
	for sizeBytes > 0 {
		writeSize := len(buffer)
		if sizeBytes < writeSize {
			writeSize = sizeBytes // write the size if the size is greater than the buffer
		}
		if _, err := file.Write(buffer[:writeSize]); err != nil {
			return err // return the error if there is an error writing the file
		}
		sizeBytes -= writeSize
	}
	return nil
}

func createMBR(mkdisk *MKDISK, sizeBytes int) error {
	var fitByte byte
	switch mkdisk.fit {
		case "FF":
			fitByte = 'F'
		case "BF":
			fitByte = 'B'
		case "WF":
			fitByte = 'W'
		default:
			fmt.Println("invalid fit")
			return errors.New("invalid fit")
	}

	// create the MBR
	mbr := &structures.MBR{
		Mbr_size:           int32(sizeBytes),
		Mbr_creation_date:  float32(time.Now().Unix()),
		Mbr_disk_signature: rand.Int31(),
		Mbr_disk_fit:       [1]byte{fitByte},
		Mbr_partitions: [4]structures.Partition{
			// initialize the partitions

			{Part_status: [1]byte{'N'}, Part_type: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'N'}, Part_correlative: -1, Part_id: [4]byte{'N'}},
			{Part_status: [1]byte{'N'}, Part_type: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'N'}, Part_correlative: -1, Part_id: [4]byte{'N'}},
			{Part_status: [1]byte{'N'}, Part_type: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'N'}, Part_correlative: -1, Part_id: [4]byte{'N'}},
			{Part_status: [1]byte{'N'}, Part_type: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'N'}, Part_correlative: -1, Part_id: [4]byte{'N'}},
		},
	}

	fmt.Println("MBR has been created")
	mbr.PrintMBR() // print the MBR

	err := mbr.SerializeMBR(mkdisk.path) // serialize the MBR
	if err != nil {
		fmt.Println("error serializing MBR:", err)
	}

	return nil
}