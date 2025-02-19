package main

import (
	"fmt"
	"runtime/debug"

	"github.com/webdestroya/aws-sso/cmd"
	"github.com/webdestroya/aws-sso/internal/utils"
)

var (
	buildVersion = "development"
	buildSha     = "unknown"
)

func main() {
	if info, ok := debug.ReadBuildInfo(); ok {
		out, _ := utils.JsonifyPretty(info)
		fmt.Printf("JSON: %s\n", out)
	}
	cmd.Execute(buildVersion, buildSha)
}
