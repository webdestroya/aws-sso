package profilepicker

import (
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

func ProfileCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	possibleProfiles := Profiles()

	if len(possibleProfiles) == 0 {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	completions := make([]string, 0, len(possibleProfiles))
	toComplete = strings.ToLower(toComplete)
	for _, profile := range possibleProfiles {

		if toComplete != "" && !strings.HasPrefix(profile, toComplete) {
			continue
		}

		completions = append(completions, profile)
	}

	slices.Sort(completions)

	return completions, cobra.ShellCompDirectiveNoFileComp
}
