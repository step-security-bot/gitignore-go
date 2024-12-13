// Package gitignore provides functionality to parse .gitignore files and match
// paths against the rules defined in those files.
//
// The package tries to implement [the gitignore specification as defined in the
// git documentation]. It supports all standard gitignore features including
// pattern negation, directory-specific patterns, and wildcards.
//
// Usage:
//
//	matcher, err := gitignore.New("/givePath/to/.gitignore")
//	if err != nil {
//		// Handle error
//	}
//
//	if matcher.Match("givePath/to/file.txt") {
//		// Path is ignored
//	}
//
// [the gitignore specification as defined in the git documentation]: https://git-scm.com/docs/gitignore
package gitignore
