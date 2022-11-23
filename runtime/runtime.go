package runtime

import (
	"gorm.io/gorm"

	"github.com/things-go/ormat/pkg/config"
)

type Runtime struct {
	Config *config.Config
	DB     *gorm.DB
}
