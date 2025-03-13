package structures

import "fmt"

type Partition struct {
	Part_status         [1]byte // N Available, 0 created, 1 mounted
	Part_type           [1]byte
	Part_fit            [1]byte
	Part_start          int32
	Part_size           int32
	Part_name           [16]byte
	Part_correlative   	int32
	Part_id 		   	[4]byte
}

func (p *Partition) CreatePartition(partStart, partSize int, partType, partFit, partName string) {

	p.Part_status[0] = '0' // 0 = Created, -1 = No created
	p.Part_start = int32(partStart)
	p.Part_size = int32(partSize)

	if len(partType) > 0 {
		p.Part_type[0] = partType[0] 
	}

	if len(partFit) > 0 {
		p.Part_fit[0] = partFit[0]
	}

	copy(p.Part_name[:], partName) // copy the name to the partition
}

// Mount a partition by id
func (p *Partition) MountPartition(correlative int, id string) error {
	p.Part_status[0] = '1' // 0 = Created, 1 = Mounted
	p.Part_correlative = int32(correlative)
	copy(p.Part_id[:], id)

	return nil
}

// GetLastEBR returns the last EBR in the linked list
func (p *Partition) GetLastEBR(path string) (*EBR, error) {
	var ebr EBR
	// get the first EBR, which is the start of the linked list
	err := ebr.DeserializeEBR(path, int64(p.Part_start))
	if err != nil {
		return nil, err
	}
	fmt.Println("First EBR")
	ebr.PrintEBR()

	// iterate over the linked list to get the last EBR
	for ebr.Ebr_part_next != -1 {
		err = ebr.DeserializeEBR(path, int64(ebr.Ebr_part_next))
		if err != nil {
			return nil, err
		}
		fmt.Println("Next EBR")
		ebr.PrintEBR()
	}

	return &ebr, nil
}

func (p *Partition) PrintPartition() {
	partStatus := rune(p.Part_status[0])
	partType := rune(p.Part_type[0])
	partFit := rune(p.Part_fit[0])

	partName := string(p.Part_name[:])
	partID := string(p.Part_id[:])

	fmt.Println("Partition")
	fmt.Printf("  Status: %c\n", partStatus)
	fmt.Printf("  Type: %c\n", partType)
	fmt.Printf("  Fit: %c\n", partFit)
	fmt.Printf("  Start: %d\n", p.Part_start)
	fmt.Printf("  Size: %d\n", p.Part_size)
	fmt.Printf("  Name: %s\n", partName)
	fmt.Printf("  Correlative: %d\n", p.Part_correlative)
	fmt.Printf("  ID: %s\n", partID)
}