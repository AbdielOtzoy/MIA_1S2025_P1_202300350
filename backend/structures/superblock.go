package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

// SuperBlock represents the superblock of a filesystem.
// It contains metadata about the filesystem, such as the number of inodes, blocks, and their sizes.
type SuperBlock struct {
	S_filesystem_type   int32   // Type of the filesystem
	S_inodes_count      int32   // Total number of inodes
	S_blocks_count      int32   // Total number of blocks
	S_free_inodes_count int32   // Number of free inodes
	S_free_blocks_count int32   // Number of free blocks
	S_mtime             float32 // Last mount time (as a Unix timestamp)
	S_umtime            float32 // Last unmount time (as a Unix timestamp)
	S_mnt_count         int32   // Number of times the filesystem has been mounted
	S_magic             int32   // Magic number to identify the filesystem
	S_inode_size        int32   // Size of each inode
	S_block_size        int32   // Size of each block
	S_first_ino         int32   // First available inode
	S_first_blo         int32   // First available block
	S_bm_inode_start    int32   // Starting position of the inode bitmap
	S_bm_block_start    int32   // Starting position of the block bitmap
	S_inode_start       int32   // Starting position of the inode table
	S_block_start       int32   // Starting position of the block table
	// Total size: 68 bytes
}

// Serialize writes the SuperBlock structure to a binary file at the specified offset.
// This function is used to persist the SuperBlock data to disk.
func (sb *SuperBlock) Serialize(path string, offset int64) error {
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

	// Serialize the SuperBlock structure directly into the file
	err = binary.Write(file, binary.LittleEndian, sb)
	if err != nil {
		return err
	}

	return nil
}

// Deserialize reads the SuperBlock structure from a binary file at the specified offset.
// This function is used to load the SuperBlock data from disk into memory.
func (sb *SuperBlock) Deserialize(path string, offset int64) error {
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

	// Get the size of the SuperBlock structure
	sbSize := binary.Size(sb)
	if sbSize <= 0 {
		return fmt.Errorf("invalid SuperBlock size: %d", sbSize)
	}

	// Read only the number of bytes corresponding to the size of the SuperBlock structure
	buffer := make([]byte, sbSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Deserialize the read bytes into the SuperBlock structure
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, sb)
	if err != nil {
		return err
	}

	return nil
}

// CreateUsersFile initializes the filesystem by creating the root directory and the users.txt file.
// This function sets up the basic structure of the filesystem.
func (sb *SuperBlock) CreateUsersFile(path string) error {
	// ----------- Create the root directory -----------
	// Create the root inode
	rootInode := &Inode{
		I_uid:   1,
		I_gid:   1,
		I_size:  0,
		I_atime: float32(time.Now().Unix()),
		I_ctime: float32(time.Now().Unix()),
		I_mtime: float32(time.Now().Unix()),
		I_block: [15]int32{sb.S_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  [1]byte{'0'},
		I_perm:  [3]byte{'7', '7', '7'},
	}

	// Serialize the root inode
	err := rootInode.Serialize(path, int64(sb.S_first_ino))
	if err != nil {
		return err
	}

	// Update the inode bitmap
	err = sb.UpdateBitmapInode(path)
	if err != nil {
		return err
	}

	// Update the superblock
	sb.S_inodes_count++
	sb.S_free_inodes_count--
	sb.S_first_ino += sb.S_inode_size

	// Create the root directory block
	rootBlock := &FolderBlock{
		B_content: [4]FolderContent{
			{B_name: [12]byte{'.'}, B_inodo: 0},
			{B_name: [12]byte{'.', '.'}, B_inodo: 0},
			{B_name: [12]byte{'-'}, B_inodo: -1},
			{B_name: [12]byte{'-'}, B_inodo: -1},
		},
	}

	// Update the block bitmap
	err = sb.UpdateBitmapBlock(path)
	if err != nil {
		return err
	}

	// Serialize the root directory block
	err = rootBlock.Serialize(path, int64(sb.S_first_blo))
	if err != nil {
		return err
	}

	// Update the superblock
	sb.S_blocks_count++
	sb.S_free_blocks_count--
	sb.S_first_blo += sb.S_block_size

	// Print the root inode
	fmt.Println("\nRoot Inode:")
	rootInode.Print()

	// Print the root directory block
	fmt.Println("\nRoot Directory Block:")
	rootBlock.Print()

	// ----------- Create /users.txt -----------
	usersText := "1,G,root\n1,U,root,123\n"

	// Deserialize the root inode
	err = rootInode.Deserialize(path, int64(sb.S_inode_start+0)) // 0 because it's the root inode
	if err != nil {
		return err
	}

	// Update the root inode's access time
	rootInode.I_atime = float32(time.Now().Unix())

	// Serialize the root inode
	err = rootInode.Serialize(path, int64(sb.S_inode_start+0)) // 0 because it's the root inode
	if err != nil {
		return err
	}

	// Deserialize the root directory block
	err = rootBlock.Deserialize(path, int64(sb.S_block_start+0)) // 0 because it's the root directory block
	if err != nil {
		return err
	}

	// Update the root directory block to include the users.txt entry
	rootBlock.B_content[2] = FolderContent{B_name: [12]byte{'u', 's', 'e', 'r', 's', '.', 't', 'x', 't'}, B_inodo: sb.S_inodes_count}

	// Serialize the root directory block
	err = rootBlock.Serialize(path, int64(sb.S_block_start+0)) // 0 because it's the root directory block
	if err != nil {
		return err
	}

	// Create the users.txt inode
	usersInode := &Inode{
		I_uid:   1,
		I_gid:   1,
		I_size:  int32(len(usersText)),
		I_atime: float32(time.Now().Unix()),
		I_ctime: float32(time.Now().Unix()),
		I_mtime: float32(time.Now().Unix()),
		I_block: [15]int32{sb.S_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  [1]byte{'1'},
		I_perm:  [3]byte{'7', '7', '7'},
	}

	// Update the inode bitmap
	err = sb.UpdateBitmapInode(path)
	if err != nil {
		return err
	}

	// Serialize the users.txt inode
	err = usersInode.Serialize(path, int64(sb.S_first_ino))
	if err != nil {
		return err
	}

	// Update the superblock
	sb.S_inodes_count++
	sb.S_free_inodes_count--
	sb.S_first_ino += sb.S_inode_size

	// Create the users.txt block
	usersBlock := &FileBlock{
		B_content: [64]byte{},
	}
	// Copy the users text into the block
	copy(usersBlock.B_content[:], usersText)

	// Serialize the users.txt block
	err = usersBlock.Serialize(path, int64(sb.S_first_blo))
	if err != nil {
		return err
	}

	// Update the block bitmap
	err = sb.UpdateBitmapBlock(path)
	if err != nil {
		return err
	}

	// Update the superblock
	sb.S_blocks_count++
	sb.S_free_blocks_count--
	sb.S_first_blo += sb.S_block_size

	// Print the updated root inode
	fmt.Println("\nUpdated Root Inode:")
	rootInode.Print()

	// Print the updated root directory block
	fmt.Println("\nUpdated Root Directory Block:")
	rootBlock.Print()

	// Print the users.txt inode
	fmt.Println("\nUsers.txt Inode:")
	usersInode.Print()

	// Print the users.txt block
	fmt.Println("\nUsers.txt Block:")
	usersBlock.Print()

	return nil
}

// PrintSuperBlock displays the values of the SuperBlock structure in a human-readable format.
// This function is used for debugging and visualization purposes.
func (sb *SuperBlock) Print() {
	// Convert mount time to a human-readable format
	mountTime := time.Unix(int64(sb.S_mtime), 0)
	// Convert unmount time to a human-readable format
	unmountTime := time.Unix(int64(sb.S_umtime), 0)

	fmt.Printf("Filesystem Type: %d\n", sb.S_filesystem_type)
	fmt.Printf("Inodes Count: %d\n", sb.S_inodes_count)
	fmt.Printf("Blocks Count: %d\n", sb.S_blocks_count)
	fmt.Printf("Free Inodes Count: %d\n", sb.S_free_inodes_count)
	fmt.Printf("Free Blocks Count: %d\n", sb.S_free_blocks_count)
	fmt.Printf("Mount Time: %s\n", mountTime.Format(time.RFC3339))
	fmt.Printf("Unmount Time: %s\n", unmountTime.Format(time.RFC3339))
	fmt.Printf("Mount Count: %d\n", sb.S_mnt_count)
	fmt.Printf("Magic: %d\n", sb.S_magic)
	fmt.Printf("Inode Size: %d\n", sb.S_inode_size)
	fmt.Printf("Block Size: %d\n", sb.S_block_size)
	fmt.Printf("First Inode: %d\n", sb.S_first_ino)
	fmt.Printf("First Block: %d\n", sb.S_first_blo)
	fmt.Printf("Bitmap Inode Start: %d\n", sb.S_bm_inode_start)
	fmt.Printf("Bitmap Block Start: %d\n", sb.S_bm_block_start)
	fmt.Printf("Inode Start: %d\n", sb.S_inode_start)
	fmt.Printf("Block Start: %d\n", sb.S_block_start)
}

// PrintInodes displays all inodes in the filesystem.
// This function is used for debugging and visualization purposes.
func (sb *SuperBlock) PrintInodes(path string) error {
	fmt.Println("\nInodes\n----------------")
	// Iterate over each inode
	for i := int32(0); i < sb.S_inodes_count; i++ {
		inode := &Inode{}
		// Deserialize the inode
		err := inode.Deserialize(path, int64(sb.S_inode_start+(i*sb.S_inode_size)))
		if err != nil {
			return err
		}
		// Print the inode
		fmt.Printf("\nInode %d:\n", i)
		inode.Print()
	}

	return nil
}

// PrintBlocks displays all blocks in the filesystem.
// This function is used for debugging and visualization purposes.
func (sb *SuperBlock) PrintBlocks(path string) error {
	fmt.Println("\nBlocks\n----------------")
	// Iterate over each inode
	for i := int32(0); i < sb.S_inodes_count; i++ {
		inode := &Inode{}
		// Deserialize the inode
		err := inode.Deserialize(path, int64(sb.S_inode_start+(i*sb.S_inode_size)))
		if err != nil {
			return err
		}
		// Iterate over each block pointer in the inode
		for _, blockIndex := range inode.I_block {
			// Skip if the block doesn't exist
			if blockIndex == -1 {
				break
			}
			// Handle folder blocks
			if inode.I_type[0] == '0' {
				block := &FolderBlock{}
				// Deserialize the folder block
				err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
				if err != nil {
					return err
				}
				// Print the folder block
				fmt.Printf("\nBlock %d:\n", blockIndex)
				block.Print()
				continue
			}
			// Handle file blocks
			if inode.I_type[0] == '1' {
				block := &FileBlock{}
				// Deserialize the file block
				err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
				if err != nil {
					return err
				}
				// Print the file block
				fmt.Printf("\nBlock %d:\n", blockIndex)
				block.Print()
				continue
			}
		}
	}

	return nil
}