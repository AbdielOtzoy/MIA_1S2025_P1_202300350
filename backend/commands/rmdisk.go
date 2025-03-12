package commands

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type RMDISK struct {
	Path string
}

func ParseRmdisk(tokens []string) (*RMDISK, error) {
	cmd := &RMDISK{} // create the rmdisk command

	args := strings.Join(tokens, " ") // join the tokens to get the arguments

	// regular expression to get the path
	re := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+`)

	matches := re.FindAllString(args, -1)

	for _, match := range matches {
		// dibide the argument in the key and value
		kv := strings.SplitN(match, "=", 3)
		if len(kv) != 2 {
			return nil, errors.New("invalid argument")
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key {
			case "-path":
				cmd.Path = value
			default:
				return nil, errors.New("invalid argument")
		}
	}

	if cmd.Path == "" {
		return nil, errors.New("missing path")
	}

	err := commandRmdisk(cmd)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func commandRmdisk(cmd *RMDISK) error {
	path := cmd.Path

	if !filepath.IsAbs(path) {
		return errors.New("invalid path")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New("disk not found")
	}

	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}