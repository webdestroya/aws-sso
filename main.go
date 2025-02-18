package main

import (
	"fmt"
	"runtime/debug"

	"github.com/webdestroya/awssso/cmd"
	"github.com/webdestroya/awssso/internal/utils"
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
