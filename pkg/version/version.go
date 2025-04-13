package version

// Version is the current version of LazyNode.
// This variable is set during build using ldflags.
var Version = "dev"

// GetVersion returns the current version string of LazyNode.
func GetVersion() string {
	return Version
}
