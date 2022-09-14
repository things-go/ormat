package view

import (
	"testing"
)

func TestName(t *testing.T) {
	t.Logf("%#v", rJSONTag.FindStringSubmatch(` 11 [@jsontag: id,omitempty,string] 11k l23123 人11`))
	t.Logf("%#v", rAffixJSONTag.FindStringSubmatch(`11大11111朋 111 11 [@affix] 11k l23123111 人11`))
}
