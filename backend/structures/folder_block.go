package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// FolderBlock represents a block in the filesystem that stores folder contents.
// It contains an array of 4 FolderContent structures, each representing a directory entry.
type FolderBlock struct {
	B_content [4]FolderContent // 4 * 16 = 64 bytes
	// Total size: 64 bytes
}

// FolderContent represents a single entry in a folder block.
// It stores the name of the entry and the inode number it points to.
type FolderContent struct {
	B_name  [12]byte // Name of the entry (up to 12 characters)
	B_inodo int32    // Inode number associated with the entry
	// Total size: 16 bytes
}

// Serialize writes the FolderBlock structure to a binary file at the specified offset.
// This function is used to persist the FolderBlock data to disk.
func (fb *FolderBlock) Serialize(path string, offset int64) error {
	// Open the file for writing or create it if it doesn't exist
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Move the file pointer to the specified offset
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Serialize the FolderBlock structure directly into the file
	err = binary.Write(file, binary.LittleEndian, fb)
	if err != nil {
		return err
	}

	return nil
}

// Deserialize reads the FolderBlock structure from a binary file at the specified offset.
// This function is used to load the FolderBlock data from disk into memory.
func (fb *FolderBlock) Deserialize(path string, offset int64) error {
	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Move the file pointer to the specified offset
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Get the size of the FolderBlock structure
	fbSize := binary.Size(fb)
	if fbSize <= 0 {
		return fmt.Errorf("invalid FolderBlock size: %d", fbSize)
	}

	// Read only the number of bytes corresponding to the size of the FolderBlock structure
	buffer := make([]byte, fbSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Deserialize the read bytes into the FolderBlock structure
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, fb)
	if err != nil {
		return err
	}

	return nil
}

// Print displays the attributes of the FolderBlock.
// This function is used for debugging and visualization purposes.
func (fb *FolderBlock) Print() {
	for i, content := range fb.B_content {
		name := string(content.B_name[:])
		fmt.Printf("Content %d:\n", i+1)
		fmt.Printf("  B_name: %s\n", name)
		fmt.Printf("  B_inodo: %d\n", content.B_inodo)
	}
}