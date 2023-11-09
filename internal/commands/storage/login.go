package storage

import (
	"fmt"
	"os/exec"
	"runtime"
)

func (h *Handler) Login(storageURL string) error {
	if err := openBrowser(storageURL); err != nil {
		h.logger.Warn(fmt.Sprintf("Unable to open browser, open the following URL: %v",
			storageURL))
	}

	return nil
}

func openBrowser(url string) error {
	var browserCommand *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		browserCommand = exec.Command("xdg-open", url)
	case "windows":
		browserCommand = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		browserCommand = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported operating system: %v", runtime.GOOS)
	}

	err := browserCommand.Run()

	return err
}
