package doctorrunner

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/factory"
)

func NewDoctorCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "doctor",
		Short:                 "Checks for possible issues using this command",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		RunE:                  RunE,
	}

	return cmd
}
