package tools

import (
	"encoding/json"
	"fmt"
	"github.com/mak-alex/backlitr/pkg/consts"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
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

func LenReadable(length int, decimals int) (out string) {
	var unit string
	var i int
	var remainder int

	// Get whole number, and the remainder for decimals
	if length > consts.TB {
		unit = "TB"
		i = length / consts.TB
		remainder = length - (i * consts.TB)
	} else if length > consts.GB {
		unit = "GB"
		i = length / consts.GB
		remainder = length - (i * consts.GB)
	} else if length > consts.MB {
		unit = "MB"
		i = length / consts.MB
		remainder = length - (i * consts.MB)
	} else if length > consts.KB {
		unit = "KB"
		i = length / consts.KB
		remainder = length - (i * consts.KB)
	} else {
		return strconv.Itoa(length) + " B"
	}

	if decimals == 0 {
		return strconv.Itoa(i) + " " + unit
	}

	// This is to calculate missing leading zeroes
	width := 0
	if remainder > consts.GB {
		width = 12
	} else if remainder > consts.MB {
		width = 9
	} else if remainder > consts.KB {
		width = 6
	} else {
		width = 3
	}

	// Insert missing leading zeroes
	remainderString := strconv.Itoa(remainder)
	for iter := len(remainderString); iter < width; iter++ {
		remainderString = "0" + remainderString
	}
	if decimals > len(remainderString) {
		decimals = len(remainderString)
	}

	return fmt.Sprintf("%d.%s %s", i, remainderString[:decimals], unit)
}
