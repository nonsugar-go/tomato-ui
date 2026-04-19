package checkpoint

import "testing"

func TestBuildKV(t *testing.T) {
	tests := []struct {
		key  string
		val  string
		want string
	}{
		{"key", "value", ` key "value"`},
	}
	for _, tt := range tests {
		got := buildKV(tt.key, tt.val)
		if got != tt.want {
			t.Errorf("buildKV(%q, %q) = %q, want %q", tt.key, tt.val, got, tt.want)
		}
	}
}

func TestBuildIndexedKV(t *testing.T) {
	tests := []struct {
		key    string
		values []string
		want   string
	}{
		{"key", []string{}, ""},
		{"key", []string{"value1"}, ` key.1 "value1"`},
		{"key", []string{"value1", "value2"}, ` key.1 "value1" key.2 "value2"`},
	}
	for _, tt := range tests {
		got := buildIndexedKV(tt.key, tt.values)
		if got != tt.want {
			t.Errorf("buildIndexedKV(%q, %v) = %q, want %q", tt.key, tt.values, got, tt.want)
		}
	}
}

func TestBuildComment(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"hello", ` comments "hello"`},
	}

	for _, tt := range tests {
		got := buildComment(tt.in)
		if got != tt.want {
			t.Errorf("input=%s got=%s want=%s", tt.in, got, tt.want)
		}
	}
}

func TestBuildTags(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want string
	}{
		{"empty", []string{}, ""},
		{"single tag", []string{"tag1"}, ` tags.1 "tag1"`},
		{"multiple tags", []string{"tag1", "tag2"}, ` tags.1 "tag1" tags.2 "tag2"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildTags(tt.in)
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}
