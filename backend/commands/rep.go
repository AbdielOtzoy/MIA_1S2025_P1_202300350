package commands

import (
	reports "backend/reports"
	stores "backend/stores"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type REP struct {
	id 				string
	path 			string
	name 			string
	path_file_ls 	string
} 

func ParseRep(tokens []string) (*REP, error) {
	cmd := &REP{}
	
	args := strings.Join(tokens, " ")

	re := regexp.MustCompile(`-id=[^\s]+|-path="[^"]+"|-path=[^\s]+|-name=[^\s]+|-path_file_ls="[^"]+"|-path_file_ls=[^\s]+`)

	matches := re.FindAllString(args, -1)

	for _, match := range matches {
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid argument: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key {
			case "-id":
				if value == "" {
					return nil, errors.New("invalid id")
				}
				cmd.id = value
			case "-path":
				if value == "" {
					return nil, errors.New("invalid path")
				}
				cmd.path = value
			case "-name":
				validNames := []string{"mbr", "disk", "inode", "block", "bm_inode", "bm_block", "sb", "file", "ls"}
				if !contains(validNames, value) {
					return nil, errors.New("nombre inválido, debe ser uno de los siguientes: mbr, disk, inode, block, bm_inode, bm_block, sb, file, ls")
				}
				cmd.name = value
			case "-path_file_ls":
				cmd.path_file_ls = value
			default:
				return nil, fmt.Errorf("invalid argument: %s", key)
		}
	}
	if cmd.id == "" || cmd.path == "" || cmd.name == "" {
		return nil, errors.New("faltan parámetros requeridos: -id, -path, -name")
	}

	err := commandRep(cmd)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return cmd, nil
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func commandRep(rep *REP) error {
	mountedMbr, mountedSb, mountedDiskPath, err := stores.GetMountedPartitionRep(rep.id)
	if err != nil {
		return err
	}

	switch rep.name {
	case "mbr":
		err = reports.ReportMBR(mountedMbr, rep.path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "inode":
		err = reports.ReportInode(mountedSb, mountedDiskPath, rep.path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "block":
		err = reports.ReportBlock(mountedSb, mountedDiskPath, rep.path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "bm_inode":
		err = reports.ReportBMInode(mountedSb, mountedDiskPath, rep.path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "bm_block":
		err = reports.ReportBMBlock(mountedSb, mountedDiskPath, rep.path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case "sb":
		err = reports.ReportSuperBlock(mountedSb, rep.path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

	}

	return nil
}