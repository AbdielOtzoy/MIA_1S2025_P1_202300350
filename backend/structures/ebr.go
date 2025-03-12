package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type EBR struct {
	Ebr_part_mount [1]byte // if the partition is mounted
	Ebr_part_fit  [1]byte // fit type: B, F, W
	Ebr_part_start int32 // starting byte of the partition
	Ebr_part_size int32 // size of the partition in bytes
	Ebr_part_next int32 // next EBR in the linked list
	Ebr_part_name [16]byte // name of the partition
}

func(ebr *EBR) SerializeEBR(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, ebr)
	if err != nil {
		return err 
	}

	return nil
}

func(ebr *EBR) DeserializeEBR(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	// get the size of the ebr struct
	ebrSize := binary.Size(ebr)
	if ebrSize <= 0 {
		return fmt.Errorf("invalid size of ebr struct: %d", ebrSize)
	}

	// read the number of bytes of the ebr struct
	buffer := make([]byte, ebrSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Deserialize the ebr struct
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, ebr)
	if err != nil {
		return err
	}

	return nil
}

func(ebr *EBR) PrintEBR() {
	partMount := rune(ebr.Ebr_part_mount[0])
	partFit := rune(ebr.Ebr_part_fit[0])
	partName := string(ebr.Ebr_part_name[:])

	fmt.Println("EBR")
	fmt.Printf("  Mount: %c\n", partMount)
	fmt.Printf("  Fit: %c\n", partFit)
	fmt.Printf("  Start: %d\n", ebr.Ebr_part_start)
	fmt.Printf("  Size: %d\n", ebr.Ebr_part_size)
	fmt.Printf("  Next: %d\n", ebr.Ebr_part_next)
	fmt.Printf("  Name: %s\n", partName)
}

func(ebr *EBR) GetFreeEBR() (*EBR, int) {
	offset := binary.Size(ebr)

	if ebr.Ebr_part_start == -1 {
		return ebr, offset
	}

	return nil, -1
}

func(ebr *EBR) GetEBRByName(name string) (*EBR, int) {
	if string(ebr.Ebr_part_name[:]) == name {
		return ebr, 0
	}

	return nil, -1
}

func(ebr *EBR) GetEBRByStart(start int32) (*EBR, int) {
	if ebr.Ebr_part_start == start {
		return ebr, 0
	}

	return nil, -1
}

func(ebr *EBR) GetEBRByNext(next int32) (*EBR, int) {
	if ebr.Ebr_part_next == next {
		return ebr, 0
	}

	return nil, -1
}


