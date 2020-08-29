package bot

import "testing"

func TestCheckURL(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"nil", "", false},
		{"tureURL", "http://baidu.com", true},
		{"char", "c", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckURL(tt.s); got != tt.want {
				t.Errorf("CheckURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
