package ssoauth

import "github.com/toqueteos/webbrowser"

func openBrowser(url string) error {
	return webbrowser.Open(url)
}
