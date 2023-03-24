package tui

import (
	"fmt"
	"github.com/pterm/pterm"
	"strings"
)

func UXError(title, msg string, err error) {
	if title != "" {
		pterm.Error.Prefix = pterm.Prefix{
			Text:  strings.ToUpper(title),
			Style: pterm.NewStyle(pterm.BgCyan, pterm.FgRed),
		}
	}

	var errMsg string
	if err != nil {
		if msg != "" {
			errMsg = fmt.Sprintf("%s: %s", msg, err)
		} else {
			errMsg = err.Error()
		}
	}

	pterm.Error.Println(errMsg)
}

func UXInfo(title, msg string) {
	if title != "" {
		pterm.Info.Prefix = pterm.Prefix{
			Text:  strings.ToUpper(title),
			Style: pterm.NewStyle(pterm.BgCyan, pterm.FgBlack),
		}
	}

	pterm.Info.Println(msg)
}

func UXSuccess(title, msg string) {
	if title != "" {
		pterm.Success.Prefix = pterm.Prefix{
			Text:  strings.ToUpper(title),
			Style: pterm.NewStyle(pterm.BgCyan, pterm.FgBlack),
		}
	}

	pterm.Success.Println(msg)
}

func UXWarning(title, msg string) {
	if title != "" {
		pterm.Warning.Prefix = pterm.Prefix{
			Text:  strings.ToUpper(title),
			Style: pterm.NewStyle(pterm.BgCyan, pterm.FgBlack),
		}
	}

	pterm.Warning.Println(msg)
}
