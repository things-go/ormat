package matcher

import (
	"testing"
)

func TestMatchEnumAnnotation(t *testing.T) {
	want := `{"0":["none","空","空注释"],"1":["key1","键1","键1注释"]}`
	tests := []struct {
		name    string
		comment string
		want    string
	}{
		{
			"",
			`11 [@enum:{"0":["none","空","空注释"],"1":["key1","键1","键1注释"]}] 11k l23123 人11`,
			want,
		},
		{
			"",
			`11 [@status:{"0":["none","空","空注释"],"1":["key1","键1","键1注释"]}] 11k l23123 人11`,
			want,
		},
		{
			"",
			`11 [@enum:  {"0":["none","空","空注释"],"1":["key1","键1","键1注释"]}  ] 11k l23123 人11`,
			want,
		},
		{
			"",
			`11, [@enum: {"0":["none","空","空注释"],"1":["key1","键1","键1注释"]}  ] [@jsontag: id,omitempty,string] 11, [@affix]`,
			want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EnumAnnotation(tt.comment); got != tt.want {
				t.Errorf("EnumAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rJSONTag(t *testing.T) {
	t.Logf("%#v", rJSONTag.FindStringSubmatch(` 11 [@jsontag: id,omitempty,string] 11k [@affix]l23123 人11`))
}

func Test_jsonTag(t *testing.T) {
	tests := []struct {
		name    string
		comment string
		want    string
	}{
		{
			"",
			"11 [@jsontag:id,omitempty,string] 11",
			"id,omitempty,string",
		},
		{
			"",
			"11 [@jsontag: id,omitempty,string] 11",
			"id,omitempty,string",
		},
		{
			"",
			"11 [@ jsontag: id,omitempty,string] 11",
			"",
		},
		{
			"",
			"11 [@jsontag: id,omitempty,string] 11, [@affix]",
			"id,omitempty,string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JsonTag(tt.comment); got != tt.want {
				t.Errorf("JsonTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rAffixJSONTag(t *testing.T) {
	t.Logf("%#v", rAffixJSONTag.FindStringSubmatch(`11大11111朋 111 11 [@affix] 11k l23123111 人11`))
}

func Test_hasAffixJSONTag(t *testing.T) {
	tests := []struct {
		name    string
		comment string
		want    bool
	}{
		{
			"",
			"朋 111 11 [@affix] 11",
			true,
		},
		{
			"",
			"朋 111 11 [@affix  ] 11",
			true,
		},
		{
			"",
			"朋 111 11 [@ affix] 11",
			false,
		},
		{
			"",
			"朋 111 11 [@affix] 11 [@jsontag:zz] xxx ",
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasAffixJSONTag(tt.comment); got != tt.want {
				t.Errorf("HasAffixJSONTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
