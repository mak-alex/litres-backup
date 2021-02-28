package tools

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func DefaultConfigPath() string {
	workDirectory := filepath.Join(func() string {
		if runtime.GOOS == "windows" {
			home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
			if home == "" {
				home = os.Getenv("USERPROFILE")
			}
			return home
		} else if runtime.GOOS == "linux" {
			home := os.Getenv("XDG_CONFIG_HOME")
			if home != "" {
				return home
			}
		}
		return os.Getenv("HOME")
	}(), ".backlitr")

	_ = MakeDirectory(workDirectory)
	_ = MakeDirectory(filepath.Join(workDirectory, "logs"))
	//_ = MakeDirectory(filepath.Join(workDirectory, "config"))
	//_ = MakeDirectory(filepath.Join(workDirectory, "config", "site"))
	//_ = MakeDirectory(filepath.Join(workDirectory, "plugins"))
	//_ = MakeDirectory(filepath.Join(workDirectory, "workspace"))
	//_ = MakeDirectory(filepath.Join(workDirectory, "workspace", "user-data"))
	//_ = MakeDirectory(filepath.Join(workDirectory, "workspace", "system-data"))
	//_ = MakeDirectory(filepath.Join(workDirectory, "workspace", "framework-temp"))

	return workDirectory
}

// MakeDirectory makes directory if is not exists
func MakeDirectory(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(dir, 0775)
		}
		return err
	}
	return nil
}

func PrettyPrint(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(b))
}
