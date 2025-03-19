package commands

import (
	stores "backend/stores"
	structures "backend/structures"
	utils "backend/utils"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type MKDIR struct {
	path string 
	p    bool  
}

/*
   mkdir -p -path=/home/user/docs/usac
   mkdir -path="/home/mis documentos/archivos clases"
*/

func ParseMkdir(tokens []string) (string, error) {
	cmd := &MKDIR{}

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`-path=[^\s]+|-p`)
	matches := re.FindAllString(args, -1)

	if len(matches) != len(tokens) {
		for _, token := range tokens {
			if !re.MatchString(token) {
				return "", fmt.Errorf("parámetro inválido: %s", token)
			}
		}
	}

	for _, match := range matches {
		kv := strings.SplitN(match, "=", 2)
		key := strings.ToLower(kv[0])

		switch key {
		case "-path":
			if len(kv) != 2 {
				return "", fmt.Errorf("formato de parámetro inválido: %s", match)
			}
			value := kv[1]
			if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
				value = strings.Trim(value, "\"")
			}
			cmd.path = value
		case "-p":
			cmd.p = true
		default:
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	if cmd.path == "" {
		return "", errors.New("faltan parámetros requeridos: -path")
	}

	err := commandMkdir(cmd)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("MKDIR: Directorio %s creado correctamente.", cmd.path), nil // Devuelve el comando MKDIR creado
}

// TODO: Aquí debería de estar logeado un usuario, por lo cual el usuario debería tener consigo el id de la partición
// En este caso el ID va a estar quemado


func commandMkdir(mkdir *MKDIR) error {
	partitionSuperblock, mountedPartition, partitionPath, err := stores.GetMountedPartitionSuperblock(idPartition)
	if err != nil {
		return fmt.Errorf("error al obtener la partición montada: %w", err)
	}

	// Crear el directorio
	err = createDirectory(mkdir.path, mkdir.p, partitionSuperblock, partitionPath, mountedPartition)
	if err != nil {
		err = fmt.Errorf("error al crear el directorio: %w", err)
	}

	return err
}

func createDirectory(dirPath string, p bool, sb *structures.SuperBlock, partitionPath string, mountedPartition *structures.Partition) error {
	fmt.Println("\nCreando directorio:", dirPath)

	parentDirs, destDir := utils.GetParentDirectories(dirPath)
	fmt.Println("\nDirectorios padres:", parentDirs)
	fmt.Println("Directorio destino:", destDir)

	if p {
		fmt.Println("Opción -p activada")
		pathSplited := utils.SplitPath(dirPath)
		for _, path := range pathSplited {
			fmt.Println("path:", path)
			// verificar que cada carpeta padre exista
			parentsDir, destDir := utils.GetParentDirectories(path)
			fmt.Printf("Directorios padres: %v\n", parentsDir)
			fmt.Printf("Directorio destino: %s\n", destDir)
			exists, err := sb.ExistsFolcer(partitionPath, parentsDir, destDir)
			if err != nil {
				return fmt.Errorf("error al verificar si la carpeta existe: %w", err)
			}

			if exists {
				fmt.Printf("La carpeta %s ya existe\n", path)
				continue
			}

			fmt.Printf("Creando la carpeta %s\n", path)
			err = sb.CreateFolder(partitionPath, parentsDir, destDir)
			if err != nil {
				return fmt.Errorf("error al crear la carpeta: %w", err)
			}

			fmt.Println("Inodos y bloques")
			sb.PrintInodes(partitionPath)
			sb.PrintBlocks(partitionPath)

			fmt.Println("Carpeta creada correctamente")
		}
	}

	// Crear el directorio segun el path proporcionado
	err := sb.CreateFolder(partitionPath, parentDirs, destDir)
	if err != nil {
		return fmt.Errorf("error al crear el directorio: %w", err)
	}

	// Imprimir inodos y bloques
	sb.PrintInodes(partitionPath)
	sb.PrintBlocks(partitionPath)

	// Serializar el superbloque
	err = sb.Serialize(partitionPath, int64(mountedPartition.Part_start))
	if err != nil {
		return fmt.Errorf("error al serializar el superbloque: %w", err)
	}

	return nil
}