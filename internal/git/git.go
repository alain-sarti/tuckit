package git

import (
	"path"
	"strings"
)

// NormalizeURL ensures the URL has a protocol.
// Short form (e.g. github.com/foo/bar) gets https:// prepended.
// SSH and full URLs are returned unchanged.
func NormalizeURL(input string) string {
	if strings.HasPrefix(input, "git@") ||
		strings.HasPrefix(input, "https://") ||
		strings.HasPrefix(input, "http://") {
		return input
	}
	return "https://" + input
}

// RepoName derives the repository name from a URL.
// Returns the last path segment with any .git suffix stripped.
func RepoName(url string) string {
	// SSH URLs: git@github.com:foo/bar.git — extract path after ":"
	if strings.HasPrefix(url, "git@") {
		parts := strings.SplitN(url, ":", 2)
		if len(parts) == 2 {
			url = parts[1]
		}
	}
	// Strip protocol prefix
	if idx := strings.Index(url, "://"); idx != -1 {
		url = url[idx+3:]
	}
	base := path.Base(url)
	return strings.TrimSuffix(base, ".git")
}
