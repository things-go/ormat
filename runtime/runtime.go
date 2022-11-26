package runtime

import (
	"gorm.io/gorm"
)

type Runtime struct {
	DB *gorm.DB
}
