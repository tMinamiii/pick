package pocket

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

// OpenBrowser open url each platform default browser.
func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		if os.Getenv("WSL_DISTRO_NAME") != "" {
			err = exec.Command("cmd.exe", "/c", "start", url).Run()
		} else {
			err = exec.Command("xdg-open", url).Run()
		}
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Run()
	case "darwin":
		err = exec.Command("open", url).Run()
	default:
	}

	if err != nil {
		log.Fatal(err)

		return
	}
}
