package reports

import (
	structures "backend/structures"
	utils "backend/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func ReportMBR(mbr *structures.MBR, diskPath string, path string) error {
	// Crear las carpetas padre si no existen
	err := utils.CreateParentDirs(path)
	if err != nil {
		return err
	}

	dotFileName, outputImage := utils.GetFileNames(path)

	dotContent := fmt.Sprintf(`digraph G {
        node [shape=plaintext]
        tabla [label=<
            <table border="0" cellborder="1" cellspacing="0">
                <tr><td colspan="2" bgcolor="blue"><font color="white"><b>REPORTE MBR</b></font></td></tr>
                <tr><td bgcolor="lightgray"><b>mbr_tamano</b></td><td>%d</td></tr>
                <tr><td bgcolor="lightgray"><b>mbr_fecha_creacion</b></td><td>%s</td></tr>
                <tr><td bgcolor="lightgray"><b>mbr_disk_signature</b></td><td>%d</td></tr>
            `, mbr.Mbr_size, time.Unix(int64(mbr.Mbr_creation_date), 0), mbr.Mbr_disk_signature)

	for i, part := range mbr.Mbr_partitions {
		partName := strings.TrimRight(string(part.Part_name[:]), "\x00")
		partStatus := rune(part.Part_status[0])
		partType := rune(part.Part_type[0])
		partFit := rune(part.Part_fit[0])

		dotContent += fmt.Sprintf(`
				<tr><td colspan="2" bgcolor="yellow"><b>PARTICIÓN %d</b></td></tr>
				<tr><td bgcolor="lightblue"><b>part_status</b></td><td>%c</td></tr>
				<tr><td bgcolor="lightblue"><b>part_type</b></td><td>%c</td></tr>
				<tr><td bgcolor="lightblue"><b>part_fit</b></td><td>%c</td></tr>
				<tr><td bgcolor="lightblue"><b>part_start</b></td><td>%d</td></tr>
				<tr><td bgcolor="lightblue"><b>part_size</b></td><td>%d</td></tr>
				<tr><td bgcolor="lightblue"><b>part_name</b></td><td>%s</td></tr>
			`, i+1, partStatus, partType, partFit, part.Part_start, part.Part_size, partName)

		if partType== 'E' {
			ebr := structures.EBR{}
			err := ebr.DeserializeEBR(diskPath, int64(part.Part_start))
			if err != nil {
				return fmt.Errorf("error al deserializar el EBR: %v", err)
			}

			ebr.PrintEBR()

			nextEBR := ebr.Ebr_part_next
			for nextEBR != -1 {
				err := ebr.DeserializeEBR(diskPath, int64(nextEBR))
				if err != nil {
					return fmt.Errorf("error al deserializar el EBR: %v", err)
				}

				ebr.PrintEBR()
				nextEBR = ebr.Ebr_part_next
			}
			
			ebrPath := utils.ReplacePath(path, "report_ebr.png")
			fmt.Println("EBR Path: ", ebrPath)

			ReportEBR(&ebr, ebrPath)
			
		}
	}

	dotContent += "</table>>] }"

	file, err := os.Create(dotFileName)
	if err != nil {
		return fmt.Errorf("error al crear el archivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(dotContent)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo: %v", err)
	}

	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", outputImage)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar el comando Graphviz: %v", err)
	}

	fmt.Println("Imagen de la tabla generada:", outputImage)
	return nil
}
