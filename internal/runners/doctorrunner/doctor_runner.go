package doctorrunner

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/appconfig"
	"github.com/webdestroya/aws-sso/internal/helpers/profilepicker"
	"github.com/webdestroya/aws-sso/internal/utils"
	"github.com/webdestroya/aws-sso/internal/utils/awsutils"
)

func RunE(cmd *cobra.Command, args []string) error {
	return newDoctor(cmd).Checkup()
}

type elDoctor struct {
	cmd *cobra.Command
}

func newDoctor(cmd *cobra.Command) *elDoctor {
	return &elDoctor{
		cmd: cmd,
	}
}

func (d elDoctor) Checkup() error {
	if ok, err := d.checkAwsConfig(); err != nil || !ok {
		return err
	}

	_ = d.checkSSOConfigs()

	return nil
}

func (d elDoctor) checkAwsConfig() (bool, error) {
	awsCfgFile := appconfig.GetAwsConfigPath()

	d.cmd.Print("Checking AWS config file...")
	_, err := os.Stat(awsCfgFile)
	if err == nil {
		d.cmd.Printf("%s (%s)\n", utils.SuccessStyle.Render("EXISTS"), awsCfgFile)
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		d.cmd.Println(utils.ErrorStyle.Render("MISSING"))
		// d.cmd.Println("Skipping configuration checks!")
		return false, nil
	}
	d.cmd.Println(utils.ErrorStyle.Render("ERROR"), err.Error())
	// d.cmd.Println("Skipping configuration checks!")
	return false, err
}

func (d elDoctor) checkSSOConfigs() error {

	d.cmd.Print("Checking for SSO profiles...")

	profiles := profilepicker.Profiles()

	if len(profiles) == 0 {
		d.cmd.Println(utils.ErrorStyle.Render("MISSING"))
		d.cmd.Println("No SSO profiles were found!")
		return nil
	}

	d.cmd.Printf("%s (%d)\n", utils.SuccessStyle.Render("FOUND"), len(profiles))

	for _, profile := range profiles {
		res := d.checkProfile(profile)

		d.cmd.Printf("  * %s %s\n", profile, res)
	}

	return nil
}

func (d elDoctor) checkProfile(profile string) string {
	cfg, err := awsutils.LoadSharedConfigProfile(d.cmd.Context(), profile)
	if err != nil {
		return utils.ErrorStyle.Render("Error:", err.Error())
	}

	if cfg.SSOStartURL != "" {
		return "(Legacy)"
	}

	return ""
}
