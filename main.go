package main

import "github.com/webdestroya/awssso/cmd"

var (
	buildVersion = "development"
	buildSha     = "unknown"
)

func main() {
	cmd.Execute(buildVersion, buildSha)
}
