package commands

import (
	"backend/stores"
	"backend/utils"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

/*
	ejemplo:
	cat -file1="path/to/file1.txt" -file2="path/to/file2.txt" -file3="path/to/file3.txt"
	parametros: -filen (n = 1, 2, 3, ..., n)
*/

type CAT struct {
	files []string
}

func ParseCat(tokens []string) (string, error) {
	cmd := &CAT{} // create the mkdisk command

	args := strings.Join(tokens, " ") // join the tokens to get the arguments

	re := regexp.MustCompile(`-file[1-9]="[^"]+"|-file[1-9]=[^\s]+`)

	matches := re.FindAllString(args, -1) // find all the matches

	for _, match := range matches {
		kv := strings.SplitN(match, "=", 2) // split the match in key and value
		if len(kv) != 2 {
			return "", fmt.Errorf("invalid argument: %s", match)
		}
		_, value := strings.ToLower(kv[0]), kv[1]

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		cmd.files = append(cmd.files, value)
	}

	if len(cmd.files) == 0 {
		return "", errors.New("missing files")
	}

	fmt.Println("CAT")
	fmt.Println(cmd.files)

	message, err := commandCat(cmd)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Se han mostrado los archivos: %s", message), nil
}

// TODO: Aquí debería de estar logeado un usuario, por lo cual el usuario debería tener consigo el id de la partición
// En este caso el ID va a estar quemado
var idPartition = "501A"

// este comando podra mostrar el contenido de un archivo
// 1. leer el superbloque de la partición
// 2. leer el inodo del archivo
// 3. leer todos los bloques del inodo y buscar el archivo
// 4. leer el contenido del archivo
// 5. mostrar el contenido del archivo
// 6. si el archivo no existe, mostrar un mensaje de error

// el cmd.files es una ruta, ir desapilando la ruta hasta llegar al archivo
// usando recursividad
func commandCat(cmd *CAT) (string, error) {
	sb, _, partitionPath, err := stores.GetMountedPartitionSuperblock(idPartition)
	if err != nil {
		return "", fmt.Errorf("error al obtener el superbloque: %w", err)
	}

	fileContents := make(map[string]string)

	for _, file := range cmd.files {
		parentDirs, destDir := utils.GetParentDirectories(file)
		fmt.Println("Parent dirs: ", parentDirs)
		fmt.Println("Dest dir: ", destDir)

		contentFile, err := sb.ReadFile(partitionPath, parentDirs, destDir)
		if err != nil {
			return "", fmt.Errorf("error al leer el archivo: %w", err)
		}

		if contentFile == "" {
			return "", fmt.Errorf("el archivo %s no existe", file)
		}

		fileContents[file] = string(contentFile)
	}

	return fmt.Sprintf("Contenido de los archivos:\n%s", fileContents), nil

}

// readFile lee el contenido de un archivo
