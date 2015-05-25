package main

// The git commit that was used to build omakase. Compiler fills it in
// with -ldflags
var GitCommit string

const Version = "0.0.1"

// This will be rendered as Version-VersionPrerelease, unless
// VersionPrerelease is empty (in which case it's a release)
const VersionPrerelease = "dev"
