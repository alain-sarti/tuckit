package git

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"github.com/foo/bar", "https://github.com/foo/bar"},
		{"https://github.com/foo/bar", "https://github.com/foo/bar"},
		{"http://github.com/foo/bar", "http://github.com/foo/bar"},
		{"git@github.com:foo/bar.git", "git@github.com:foo/bar.git"},
	}
	for _, tc := range cases {
		got := NormalizeURL(tc.input)
		if got != tc.want {
			t.Errorf("NormalizeURL(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestRepoName(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"https://github.com/foo/bar", "bar"},
		{"https://github.com/foo/bar.git", "bar"},
		{"git@github.com:foo/bar.git", "bar"},
		{"git@github.com:bar.git", "bar"}, // single-level SSH path
		{"github.com/foo/bar", "bar"},
	}
	for _, tc := range cases {
		got := RepoName(tc.input)
		if got != tc.want {
			t.Errorf("RepoName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
