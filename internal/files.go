package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path"
)

func (c *Config) GetAllFilesData() error {
	root := os.DirFS("networks/")

	txtFiles, err := fs.Glob(root, "*.txt")
	if err != nil {
		return err
	}
	fmt.Println(txtFiles)

	c.FileNames = make([]string, 0)
	c.FileStations = make(map[string][]string, 0)

	for _, v := range txtFiles {
		c.FileNames = append(c.FileNames, path.Join("networks/", v))
	}

	c.FileStations = make(map[string][]string, 0)
	for _, name := range c.FileNames {
		c.FileStations[name] = make([]string, 0)
	}

	for _, name := range c.FileNames {
		c.NetworkPath = name
		c.ParseAndValidateNetwork()
	}

	return nil
}
