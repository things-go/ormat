package command

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/things-go/ens/driver"
)

func LoadDriver(URL string) (driver.Driver, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	d, ok := driver.LoadDriver(u.Scheme)
	if !ok {
		return nil, fmt.Errorf("unsupported schema, only support [%v]", strings.Join(driver.DriverNames(), ", "))
	}
	return d, nil
}

func joinFilename(dir, filename, suffix string) string {
	suffix = strings.TrimSpace(suffix)
	if suffix != "" && !strings.HasPrefix(suffix, ".") {
		suffix = "." + suffix
	}
	return filepath.Join(dir, filename) + suffix
}

// WriteFile writes data to a file named by filename.
// If the file does not exist, WriteFile creates it
// and its upper level paths.
func WriteFile(filename string, data []byte) error {
	if err := os.MkdirAll(path.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0655)
}
