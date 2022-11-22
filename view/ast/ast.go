package ast

// interval
const delimTab = "\t"
const delimLF = "\n"
const delimSpace = " "

// ImportsHeads import head options
var ImportsHeads = map[string]string{
	"string":         `"string"`,
	"time.Time":      `"time"`,
	"gorm.Model":     `"gorm.io/gorm"`,
	"fmt":            `"fmt"`,
	"datatypes.JSON": `"gorm.io/datatypes"`,
	"datatypes.Date": `"gorm.io/datatypes"`,
}
