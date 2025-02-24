package doctorrunner

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/appconfig"
	"github.com/webdestroya/aws-sso/internal/runners/listrunner"
	"github.com/webdestroya/aws-sso/internal/utils"
)

func RunE(cmd *cobra.Command, args []string) error {

	// TODO: check for aws cli?
	return newDoctor(cmd).Checkup()

	awsCfgFile := appconfig.GetAwsConfigPath()

	cmd.Printf("Checking %s file...", awsCfgFile)
	if _, err := os.Stat(awsCfgFile); err == nil {
		cmd.Printf("EXISTS (%s)\n", awsCfgFile)
	} else if errors.Is(err, os.ErrNotExist) {
		cmd.Println(utils.ErrorStyle.Render("MISSING"))
		cmd.Println("Skipping configuration checks!")
		return nil
	} else {
		cmd.Println(utils.ErrorStyle.Render("ERROR"), err.Error())
		cmd.Println("Skipping configuration checks!")
		return nil
	}

	cmd.Print("Checking for sso configurations...")

	entries, err := listrunner.GetSSOEntries()
	if err != nil {
		return err
	}
	cmd.Printf("FOUND (%d)\n", len(entries))

	for _, entry := range entries {
		cmd.Printf(" * %s ", entry.String())

		if len(entry.Profiles) > 0 {
			cmd.Printf("(Used By: %s)", strings.Join(entry.Profiles, ", "))
		} else {
			cmd.Print("(Not used by any profiles)")
		}
		cmd.Println()
	}

	if len(entries) == 0 {
		cmd.Println(utils.ErrorStyle.Render("NONE"), "No 'sso-session' entries found. You need to configure SSO!")
		return nil
	}

	return nil
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
	return nil
}
