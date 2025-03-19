package reports

import (
	structures "backend/structures"
	utils "backend/utils"
	"os"
	"os/exec"
)

func ReportBlock(superblock *structures.SuperBlock, diskPath string, path string) error {
	err := utils.CreateParentDirs(path)
	if err != nil {
		return err
	}

	dotFileName, outputImage := utils.GetFileNames(path)

	dotContent := `digraph G {
		rankdir=LR; // Alineación horizontal
		node [shape=plaintext]
	`

	// Iterar sobre cada bloque
	for i := int32(0); i < superblock.S_blocks_count; i++ {
		
		
	}

	dotContent += "}"

	// Crear archivo DOT
	file, err := os.Create(dotFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(dotContent)
	if err != nil {
		return err
	}

	// Generar imagen con Graphviz
	cmd := exec.Command("dot", "-T"+outputImage, dotFileName, "-o", path)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}