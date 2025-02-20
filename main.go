package main

import (
	"github.com/webdestroya/aws-sso/cmd"
)

var (
	buildVersion = "development"
	buildSha     = "unknown"
)

func main() {
	cmd.Execute(buildVersion, buildSha)
}
