package utils

import (
	"reflect"
	"testing"
)

func Test_getTLDMap(t *testing.T) {
	type args struct {
		domains []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]struct{}
		wantErr bool
	}{
		{"case1", args{
			[]string{"example.com", "www.example.com", "test.example.com", "a.b.c.example.com"},
		}, map[string]struct{}{"example.com": {}}, false},
		{"edu.cn", args{
			[]string{"example.edu.cn", "www.example.edu.cn", "test.example.edu.cn", "a.b.c.example.edu.cn", "www.example.cn"},
		}, map[string]struct{}{"example.edu.cn": {}, "example.cn": {}}, false},
		{"com.cn", args{
			[]string{"example.com.cn", "www.example.com.cn", "test.example.com.cn", "a.b.c.example.com.cn"},
		}, map[string]struct{}{"example.com.cn": {}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTLDMap(tt.args.domains)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTLDMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTLDMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}
