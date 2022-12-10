package ast

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
			if got := MatchEnumAnnotation(tt.comment); got != tt.want {
				t.Errorf("MatchEnumAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}
