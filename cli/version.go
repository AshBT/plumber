package cli

// The git commit that was used to build plumb. Compiler fills it in
// with -ldflags
var GitCommit string

const version = "0.0.1"

// This will be rendered as Version-VersionPrerelease, unless
// VersionPrerelease is empty (in which case it's a release)
const versionPrerelease = "dev"
