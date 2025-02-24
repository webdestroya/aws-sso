package main

import (
	"os"

	"github.com/webdestroya/aws-sso/cmd"
)

var (
	buildVersion = "development"
	buildSha     = "unknown"
)

func main() {
	code := cmd.Execute(buildVersion, buildSha)
	os.Exit(code)
}
