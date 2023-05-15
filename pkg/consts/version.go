package consts

import (
	"fmt"
	"runtime"
)

const Version = "v0.9.2"

func BuildVersion() string {
	return fmt.Sprintf("%s\nGo Version: %s\nGo Os: %s\nGo Arch: %s\n",
		Version, runtime.Version(),
		runtime.GOOS, runtime.GOARCH)
}
