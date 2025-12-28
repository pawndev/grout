package utils

import gaba "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool"

var defaultStatusBar = gaba.StatusBarOptions{
	Enabled:  true,
	ShowTime: false,
	Icons:    []gaba.StatusBarIcon{},
}

func StatusBar() gaba.StatusBarOptions {
	return defaultStatusBar
}

func AddIcon(icon gaba.StatusBarIcon) {
	defaultStatusBar.Icons = append(defaultStatusBar.Icons, icon)
}
