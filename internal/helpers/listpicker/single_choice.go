package listpicker

import (
	"slices"

	"github.com/charmbracelet/huh"
)

func NewSingleChoice(title string, choices []string) (string, error) {

	choiceList := make([]huh.Option[string], 0, len(choices))
	slices.Sort(choices)
	for _, c := range choices {
		choiceList = append(choiceList, huh.NewOption(c, c))
	}

	choice := ""

	f := huh.NewSelect[string]().
		Title(title).
		Options(choiceList...).
		Value(&choice).
		Description(" ").
		WithTheme(huh.ThemeCharm())

	if err := f.Run(); err != nil {
		// huh.ErrUserAborted
		return "", err
	}

	return choice, nil
}
