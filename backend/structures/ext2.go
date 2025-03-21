package structures

import (
	utils "backend/utils"
	"strings"
	"time"
	"fmt"
)

// Crear users.txt en nuestro sistema de archivos
func (sb *SuperBlock) CreateUsersFile(path string) error {
	// ----------- Creamos / -----------
	// Creamos el inodo raíz
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

	// Serializar el inodo raíz
	err := rootInode.Serialize(path, int64(sb.S_first_ino))
	if err != nil {
		return err
	}

	// Actualizar el bitmap de inodos
	err = sb.UpdateBitmapInode(path)
	if err != nil {
		return err
	}

	// Actualizar el superbloque
	sb.S_inodes_count++
	sb.S_free_inodes_count--
	sb.S_first_ino += sb.S_inode_size

	// Creamos el bloque del Inodo Raíz
	rootBlock := &FolderBlock{
		B_content: [4]FolderContent{
			{B_name: [12]byte{'.'}, B_inodo: 0},
			{B_name: [12]byte{'.', '.'}, B_inodo: 0},
			{B_name: [12]byte{'-'}, B_inodo: -1},
			{B_name: [12]byte{'-'}, B_inodo: -1},
		},
	}

	// Actualizar el bitmap de bloques
	err = sb.UpdateBitmapBlock(path)
	if err != nil {
		return err
	}

	// Serializar el bloque de carpeta raíz
	err = rootBlock.Serialize(path, int64(sb.S_first_blo))
	if err != nil {
		return err
	}

	// Actualizar el superbloque
	sb.S_blocks_count++
	sb.S_free_blocks_count--
	sb.S_first_blo += sb.S_block_size

	// ----------- Creamos /users.txt -----------
	usersText := "1,G,root\n1,U,root,root,123\n"

	// Deserializar el inodo raíz
	err = rootInode.Deserialize(path, int64(sb.S_inode_start+0)) // 0 porque es el inodo raíz
	if err != nil {
		return err
	}

	// Actualizamos el inodo raíz
	rootInode.I_atime = float32(time.Now().Unix())

	// Serializar el inodo raíz
	err = rootInode.Serialize(path, int64(sb.S_inode_start+0)) // 0 porque es el inodo raíz
	if err != nil {
		return err
	}

	// Deserializar el bloque de carpeta raíz
	err = rootBlock.Deserialize(path, int64(sb.S_block_start+0)) // 0 porque es el bloque de carpeta raíz
	if err != nil {
		return err
	}

	// Actualizamos el bloque de carpeta raíz
	rootBlock.B_content[2] = FolderContent{B_name: [12]byte{'u', 's', 'e', 'r', 's', '.', 't', 'x', 't'}, B_inodo: sb.S_inodes_count}

	// Serializar el bloque de carpeta raíz
	err = rootBlock.Serialize(path, int64(sb.S_block_start+0)) // 0 porque es el bloque de carpeta raíz
	if err != nil {
		return err
	}

	// Creamos el inodo users.txt
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

	// Actualizar el bitmap de inodos
	err = sb.UpdateBitmapInode(path)
	if err != nil {
		return err
	}

	// Serializar el inodo users.txt
	err = usersInode.Serialize(path, int64(sb.S_first_ino))
	if err != nil {
		return err
	}

	// Actualizamos el superbloque
	sb.S_inodes_count++
	sb.S_free_inodes_count--
	sb.S_first_ino += sb.S_inode_size

	// Creamos el bloque de users.txt
	usersBlock := &FileBlock{
		B_content: [64]byte{},
	}
	// Copiamos el texto de usuarios en el bloque
	copy(usersBlock.B_content[:], usersText)

	// Serializar el bloque de users.txt
	err = usersBlock.Serialize(path, int64(sb.S_first_blo))
	if err != nil {
		return err
	}

	// Actualizar el bitmap de bloques
	err = sb.UpdateBitmapBlock(path)
	if err != nil {
		return err
	}

	// Actualizamos el superbloque
	sb.S_blocks_count++
	sb.S_free_blocks_count--
	sb.S_first_blo += sb.S_block_size

	return nil
}

// createFolderInInode crea una carpeta en un inodo específico
func (sb *SuperBlock) createFolderInInode(path string, inodeIndex int32, parentsDir []string, destDir string) error {
	// Crear un nuevo inodo
	inode := &Inode{}
	// Deserializar el inodo
	err := inode.Deserialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
	if err != nil {
		return err
	}
	// Verificar si el inodo es de tipo carpeta
	if inode.I_type[0] == '1' {
		return nil
	}

	// Iterar sobre cada bloque del inodo (apuntadores)
	for _, blockIndex := range inode.I_block {
		// Si el bloque no existe, salir
		if blockIndex == -1 {
			break
		}

		// Crear un nuevo bloque de carpeta
		block := &FolderBlock{}

		block.Print()

		// Deserializar el bloque
		err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size))) // 64 porque es el tamaño de un bloque
		if err != nil {
			return err
		}

		// Iterar sobre cada contenido del bloque, desde el index 2 porque los primeros dos son . y ..
		for indexContent := 2; indexContent < len(block.B_content); indexContent++ {
			// Obtener el contenido del bloque
			content := block.B_content[indexContent]

			// Sí las carpetas padre no están vacías debereamos buscar la carpeta padre más cercana
			if len(parentsDir) != 0 {
				//fmt.Println("---------ESTOY  VISITANDO--------")

				// Si el contenido está vacío, salir
				if content.B_inodo == -1 {
					break
				}

				// Obtenemos la carpeta padre más cercana
				parentDir, err := utils.First(parentsDir)
				if err != nil {
					return err
				}

				// Convertir B_name a string y eliminar los caracteres nulos
				contentName := strings.Trim(string(content.B_name[:]), "\x00 ")
				// Convertir parentDir a string y eliminar los caracteres nulos
				parentDirName := strings.Trim(parentDir, "\x00 ")
				// Si el nombre del contenido coincide con el nombre de la carpeta padre
				if strings.EqualFold(contentName, parentDirName) {
					//fmt.Println("---------LA ENCONTRÉ-------")
					// Si son las mismas, entonces entramos al inodo que apunta el bloque
					err := sb.createFolderInInode(path, content.B_inodo, utils.RemoveElement(parentsDir, 0), destDir)
					if err != nil {
						return err
					}
					return nil
				}
			} else {
				//fmt.Println("---------ESTOY  CREANDO--------")

				// Si el apuntador al inodo está ocupado, continuar con el siguiente
				if content.B_inodo != -1 {
					continue
				}

				// Actualizar el contenido del bloque
				copy(content.B_name[:], destDir)
				content.B_inodo = sb.S_inodes_count

				// Actualizar el bloque
				block.B_content[indexContent] = content

				// Serializar el bloque
				err = block.Serialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
				if err != nil {
					return err
				}

				// Crear el inodo de la carpeta
				folderInode := &Inode{
					I_uid:   1,
					I_gid:   1,
					I_size:  0,
					I_atime: float32(time.Now().Unix()),
					I_ctime: float32(time.Now().Unix()),
					I_mtime: float32(time.Now().Unix()),
					I_block: [15]int32{sb.S_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
					I_type:  [1]byte{'0'},
					I_perm:  [3]byte{'6', '6', '4'},
				}

				// Serializar el inodo de la carpeta
				err = folderInode.Serialize(path, int64(sb.S_first_ino))
				if err != nil {
					return err
				}

				// Actualizar el bitmap de inodos
				err = sb.UpdateBitmapInode(path)
				if err != nil {
					return err
				}

				// Actualizar el superbloque
				sb.S_inodes_count++
				sb.S_free_inodes_count--
				sb.S_first_ino += sb.S_inode_size

				// Crear el bloque de la carpeta
				folderBlock := &FolderBlock{
					B_content: [4]FolderContent{
						{B_name: [12]byte{'.'}, B_inodo: content.B_inodo},
						{B_name: [12]byte{'.', '.'}, B_inodo: inodeIndex},
						{B_name: [12]byte{'-'}, B_inodo: -1},
						{B_name: [12]byte{'-'}, B_inodo: -1},
					},
				}

				// Serializar el bloque de la carpeta
				err = folderBlock.Serialize(path, int64(sb.S_first_blo))
				if err != nil {
					return err
				}

				// Actualizar el bitmap de bloques
				err = sb.UpdateBitmapBlock(path)
				if err != nil {
					return err
				}

				// Actualizar el superbloque
				sb.S_blocks_count++
				sb.S_free_blocks_count--
				sb.S_first_blo += sb.S_block_size

				return nil
			}
		}

	}
	return nil
}

func (sb *SuperBlock) createFileInode(path string, inodeIndex int32, parentsDir []string, destDir string, r bool, size int, contentFile string) error {
	// crear un nuevo inodo
	inode := &Inode{}
	// deserializar el inodo
	err := inode.Deserialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
	if err != nil {
		return err
	}

	// verificar que el inodo sea de tipo carpeta
	if inode.I_type[0] == '1' {
		return nil
	}

	// Iterar sobre cada bloque del inodo (apuntadores)
	for _, blockIndex := range inode.I_block {
		if blockIndex == -1 {
			break
		}

		block := &FolderBlock{}
		// deserializar el bloque
		err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size))) // 64 porque es el tamaño de un bloque
		if err != nil {
			return err
		}

		// Iterar sobre cada contenido del bloque, desde el index 2 porque los primeros dos son . y ..
		for indexContent := 2; indexContent < len(block.B_content); indexContent++ {
			content := block.B_content[indexContent]

			// Si las carpetas padre no están vacías debereamos buscar la carpeta padre más cercana
			if len(parentsDir) != 0 {
				// Si el contenido está vacío, salir
				if content.B_inodo == -1 {
					break
				}

				// Obtenemos la carpeta padre más cercana
				parentDir, err := utils.First(parentsDir)
				if err != nil {
					return err
				}

				// Convertir B_name a string y eliminar los caracteres nulos
				contentName := strings.Trim(string(content.B_name[:]), "\x00 ")
				// Convertir parentDir a string y eliminar los caracteres nulos
				parentDirName := strings.Trim(parentDir, "\x00 ")
				// Si el nombre del contenido coincide con el nombre de la carpeta padre
				if strings.EqualFold(contentName, parentDirName) {
					err := sb.createFileInode(path, content.B_inodo, utils.RemoveElement(parentsDir, 0), destDir, r, size, contentFile)
					if err != nil {
						return err
					}
					return nil
				}
			} else {
				fmt.Println("---------ESTOY  CREANDO--------")
				// Si el apuntador al inodo está ocupado, continuar con el siguiente
				if content.B_inodo != -1 {
					continue
				}

				// Actualizar el contenido del bloque
				copy(content.B_name[:], destDir)
				content.B_inodo = sb.S_inodes_count
				// Actualizar el bloque
				block.B_content[indexContent] = content
				// Serializar el bloque
				err = block.Serialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size)))
				if err != nil {
					return err
				}

				// Crear el inodo del archivo
				// Crear el inodo de la carpeta
				fileInode := &Inode{
					I_uid:   1,
					I_gid:   1,
					I_size:  int32(size),
					I_atime: float32(time.Now().Unix()),
					I_ctime: float32(time.Now().Unix()),
					I_mtime: float32(time.Now().Unix()),
					I_block: [15]int32{sb.S_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
					I_type:  [1]byte{'1'},
					I_perm:  [3]byte{'6', '6', '4'},
				}

				// Serializar el inodo del archivo
				err = fileInode.Serialize(path, int64(sb.S_first_ino))
				if err != nil {
					return err
				}

				// Actualizar el bitmap de inodos
				err = sb.UpdateBitmapInode(path)
				if err != nil {
					return err
				}

				// Actualizar el superbloque
				sb.S_inodes_count++
				sb.S_free_inodes_count--
				sb.S_first_ino += sb.S_inode_size
				// Crear el bloque del archivo
				// colocar contenido de prueba

				// convertir el contentFile a 64 bytes
				contentByte := [64]byte{}
				copy(contentByte[:], contentFile)
				fileBlock := &FileBlock{
					B_content: contentByte,
				}

				fmt.Println("Bloque de archivo: " )
				fileBlock.Print()

				// Serializar el bloque del archivo
				err = fileBlock.Serialize(path, int64(sb.S_first_blo))
				if err != nil {
					return err
				}

				// Actualizar el bitmap de bloques
				err = sb.UpdateBitmapBlock(path)
				if err != nil {
					return err
				}

				// Actualizar el superbloque
				sb.S_blocks_count++
				sb.S_free_blocks_count--
				sb.S_first_blo += sb.S_block_size
			
				return nil
			
			}
		}
	}
	return nil
}

func (sb *SuperBlock) folderExists(path string, inodeIndex int32, parentsDir []string, destDir string) (bool, error) {
	// Crear un nuevo inodo
	inode := &Inode{}
	// Deserializar el inodo
	err := inode.Deserialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
	if err != nil {
		return false, err
	}
	// Verificar si el inodo es de tipo carpeta
	if inode.I_type[0] == '1' {
		return false, nil
	}

	// Iterar sobre cada bloque del inodo (apuntadores)
	for _, blockIndex := range inode.I_block {
		// Si el bloque no existe, salir
		if blockIndex == -1 {
			break
		}

		// Crear un nuevo bloque de carpeta
		block := &FolderBlock{}
		
		// Deserializar el bloque
		err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size))) // 64 porque es el tamaño de un bloque
		if err != nil {
			return false, err
		}
		block.Print()

		// Iterar sobre cada contenido del bloque, desde el index 2 porque los primeros dos son . y ..
		for indexContent := 2; indexContent < len(block.B_content); indexContent++ {
			// Obtener el contenido del bloque
			content := block.B_content[indexContent]

			if len(parentsDir) != 0 {
				if content.B_inodo == -1 {
					break
				}
				// Obtenemos la carpeta padre más cercana
				parentDir, err := utils.First(parentsDir)
				if err != nil {
					return false, err
				}

				// Convertir B_name a string y eliminar los caracteres nulos
				contentName := strings.Trim(string(content.B_name[:]), "\x00 ")
				// Convertir parentDir a string y eliminar los caracteres nulos
				parentDirName := strings.Trim(parentDir, "\x00 ")

				// Si el nombre del contenido coincide con el nombre de la carpeta padre
				if strings.EqualFold(contentName, parentDirName) {
					fmt.Println("entrando a la carpeta padre")
					fileContent, err := sb.folderExists(path, content.B_inodo, utils.RemoveElement(parentsDir, 0), destDir)
					if err != nil {
						return false, err
					}
					return fileContent, nil
				}
			} else {
				
				if content.B_inodo == -1 {
					continue
				}
				fmt.Println("ya no hay más carpetas padre")
				fmt.Println("Carpeta destino: ", destDir)

				// convertir content.B_name a string
				contentName := strings.Trim(string(content.B_name[:]), "\x00 ")
				fmt.Println("Nombre de la carpeta: ", contentName)
				if(contentName == destDir) {
					fmt.Println("---------LA ENCONTRÉ-------")
					return true, nil
				}
			}
		}

	}
	return false, nil
}


func (sb *SuperBlock) readFileInInode(path string, inodeIndex int32, parentsDir []string, destDir string) (string, error) {
	// Crear un nuevo inodo
	inode := &Inode{}
	// Deserializar el inodo
	err := inode.Deserialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
	if err != nil {
		return "", err
	}
	// Verificar si el inodo es de tipo carpeta
	if inode.I_type[0] == '1' {
		return "", nil
	}

	fmt.Println("destDir: ", destDir)
	fmt.Println("Padres: ", parentsDir)

	fmt.Println("Inodo: ", inodeIndex)


	// Iterar sobre cada bloque del inodo (apuntadores)
	for _, blockIndex := range inode.I_block {
		// Si el bloque no existe, salir
		if blockIndex == -1 {
			break
		}

		// Crear un nuevo bloque de carpeta
		block := &FolderBlock{}

		// Deserializar el bloque
		err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size))) // 64 porque es el tamaño de un bloque
		if err != nil {
			return "", err
		}


		fmt.Println("Bloque: ", blockIndex)
		block.Print()
		
		// Iterar sobre cada contenido del bloque, desde el index 2 porque los primeros dos son . y ..
		for indexContent := 2; indexContent < len(block.B_content); indexContent++ {
			// Obtener el contenido del bloque
			content := block.B_content[indexContent]

			if len(parentsDir) != 0 {
				fmt.Println("---------ESTOY  VISITANDO--------")

				if content.B_inodo == -1 {
					break
				}
				// Obtenemos la carpeta padre más cercana
				parentDir, err := utils.First(parentsDir)
				if err != nil {
					return "", err
				}

				// Convertir B_name a string y eliminar los caracteres nulos
				contentName := strings.Trim(string(content.B_name[:]), "\x00 ")
				// Convertir parentDir a string y eliminar los caracteres nulos
				parentDirName := strings.Trim(parentDir, "\x00 ")
				
				fmt.Println("Nombre del contenido: ", contentName)
				fmt.Println("Nombre del archivo: ", destDir)
				fmt.Println("Nombre de la carpeta padre: ", parentDirName)

				// Si el nombre del contenido coincide con el nombre de la carpeta padre
				if strings.EqualFold(contentName, parentDirName) {
					fmt.Println("entrando a la carpeta padre")
					fileContent, err := sb.readFileInInode(path, content.B_inodo, utils.RemoveElement(parentsDir, 0), destDir)
					if err != nil {
						return "", err
					}
					if fileContent != "" {
						return fileContent, nil
					}
					return "", fmt.Errorf("no se encontró el archivo")
				}
			} else {
				fmt.Println("ya no hay más carpetas padre")
				if content.B_inodo == -1 {
					continue
				}
				// convertir destDir a 12 bytes
				destDirByte := [12]byte{}
				copy(destDirByte[:], destDir)

				if(content.B_name == destDirByte) {
					fmt.Println("---------LA ENCONTRÉ-------")
					// Si son las mismas, entonces entramos al inodo que apunta el bloque
					inodeFile := &Inode{}
					err := inodeFile.Deserialize(path, int64(sb.S_inode_start+(content.B_inodo*sb.S_inode_size)))
					if err != nil {
						return "", err
					}

					// Verificar si el inodo es de tipo archivo
					if inodeFile.I_type[0] == '1' {
						// Crear un nuevo bloque de archivo
						blockFile := &FileBlock{}

						// Deserializar el bloque de archivo
						err := blockFile.Deserialize(path, int64(sb.S_block_start+(inodeFile.I_block[0]*sb.S_block_size)))
						if err != nil {
							return "", err
						}

						// Imprimir el contenido del bloque de archivo
						fmt.Println("Contenido del archivo: ")
						fmt.Println(string(blockFile.B_content[:]))

						return string(blockFile.B_content[:]), nil
						
					} 
					return "", fmt.Errorf("el inodo no es de tipo archivo")
				}
			}
		}
	}
	return "", fmt.Errorf("no se encontró el archivo")
}   

func (sb *SuperBlock) loginUserInInode(user string, password string, path string) error {
	// para iniciar sesión se debe leer el archivo users.txt
	// y verificar si el usuario y la contraseña son correctos

	// el users.txt siempre esta en el bloque 1
	useBlock := &FileBlock{}

	// deserializar el bloque
	err := useBlock.Deserialize(path, int64(sb.S_block_start+(1*sb.S_block_size)))
	if err != nil {
		return err
	}

	// obtener el contenido del bloque
	content := string(useBlock.B_content[:])
	fmt.Println("Contenido del bloque: ", content)

	// separar el contenido por lineas
	lines := strings.Split(content, "\n")

	// iterar sobre cada linea
	for _, line := range lines {
		fmt.Println("Linea: ", line)
		// separar la linea por comas
		parts := strings.Split(line, ",")

		/*
			formato del archivo:
			GID, Tipo, Grupo \n
			UID, Tipo, Grupo, Usuario, Contraseña \n
		*/

		fmt.Println("Partes: ", parts)
		

		// verificar si es un grupo
		if parts[1] == "G" {
			fmt.Println("Grupo: ", parts[2])
		}
		// verificar si es un usuario
		if parts[1] == "U" {
			fmt.Println("Usuario: ", parts[3])
			fmt.Println("Contraseña: ", parts[4])
			if parts[3] == user && parts[4] == password {
				fmt.Println("Usuario: ", parts[3])
				fmt.Println("Contraseña: ", parts[4])
				return nil
			} else if parts[3] == user && parts[4] != password {
				fmt.Println("Contraseña incorrecta")
				return fmt.Errorf("contraseña incorrecta")
			} else {
				fmt.Println("Usuario o contraseña incorrectos")
			}
		}



	}

	return fmt.Errorf("no se encontró el usuario")

}
