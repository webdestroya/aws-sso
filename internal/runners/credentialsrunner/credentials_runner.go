package credentialsrunner

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/helpers/getcreds"
	"github.com/webdestroya/aws-sso/internal/utils"
)

func RunE(cmd *cobra.Command, args []string) error {

	profile := args[0]

	credentials, err := getcreds.GetAWSCredentials(cmd.Context(), cmd.ErrOrStderr(), profile, getcreds.WithLoginDisabled())
	if err != nil {
		return err
	}

	creds := map[string]any{
		"Version":         1,
		"AccessKeyId":     credentials.AccessKeyID,
		"SecretAccessKey": credentials.SecretAccessKey,
		"SessionToken":    credentials.SessionToken,
		"Expiration":      credentials.Expires.Format(time.RFC3339),
	}

	pretty, err := utils.JsonifyPretty(creds)
	if err != nil {
		return err
	}

	cmd.Println(pretty)

	return nil
}
