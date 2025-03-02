package cmd

import "github.com/webdestroya/aws-sso/internal/runners/doctorrunner"

func init() {
	rootCmd.AddCommand(doctorrunner.NewDoctorCmd(cmdFactory))
}
