package tui

import (
	"github.com/pterm/pterm"
	"strings"
)

type TUITitle struct {
}

func (t *TUITitle) ShowTitleAndDescription(title, description string) {
	titleNormalised := strings.TrimSpace(strings.ToUpper(title))
	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString(titleNormalised)).
		Srender()
	pterm.DefaultCenter.Println(s)

	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(description)
}

func (t *TUITitle) ShowTitle(title string) {
	titleNormalised := strings.TrimSpace(strings.ToUpper(title))
	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString(titleNormalised)).
		Srender()
	pterm.DefaultCenter.Println(s)
}

func (t *TUITitle) ShowSubTitle(subtitle string) {
	subtitleNormalised := strings.TrimSpace(strings.ToUpper(subtitle))
	pterm.Println()
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println("--------------------------------")
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(subtitleNormalised)
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println("--------------------------------")
	pterm.Println()
}

func NewTitle() TUIDisplayer {
	return &TUITitle{}
}
