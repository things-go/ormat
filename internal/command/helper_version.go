package command

import (
	"fmt"
	"runtime"
)

const version = "v0.13.1-rc2"

func BuildVersion() string {
	return fmt.Sprintf("%s\nGo Version: %s\nGo Os: %s\nGo Arch: %s\n",
		version, runtime.Version(),
		runtime.GOOS, runtime.GOARCH)
}
