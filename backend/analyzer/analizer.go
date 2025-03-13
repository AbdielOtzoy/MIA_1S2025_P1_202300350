package analyzer

import (
	commands "backend/commands"
	"errors"
	"fmt"
	"strings"
)

func Analyzer(input string) (interface{}, error) {
	
	tokens := strings.Fields(input)

	if len(tokens) == 0 {
		return nil, errors.New("no se proporcionó ningún comando")
	}

	switch tokens[0] {
	case "mkdir":
		return ParseMkdir(tokens[1:])
	case "mkdisk":
		return commands.ParseMkdisk(tokens[1:])
	case "rmdisk":
		return commands.ParseRmdisk(tokens[1:])
	case "fdisk":
		return commands.ParseFdisk(tokens[1:])
	case "mount":
		return commands.ParseMount(tokens[1:])
	case "mkfs":
		return commands.ParseMkfs(tokens[1:])
	case "rep":
		return commands.ParseRep(tokens[1:])
	case "mounted":
		return commands.ParseMounted(tokens[1:])
	default:
		return nil, fmt.Errorf("comando desconocido: %s", tokens[0])
	}
}