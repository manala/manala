package web

import (
	"fmt"
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) {
	switch runtime.GOOS {
	case "darwin":
		//err = exec.Command("open", url).Start()
		if err := exec.Command("open", url).Start(); err != nil {
			// ...
			fmt.Println(err)
		}
		//case "linux":
		//	err = exec.Command("xdg-open", url).Start()
		//case "windows":
		//	err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		//default:
		//	err = fmt.Errorf("unsupported platform")
		//}
		//if err != nil {
		//	log.Fatal(err)
		//}
	}
}

/*
case "windows":
cmd = "cmd"
args = []string{"/c", "start"}
case "darwin":
cmd = "open"
default: // "linux", "freebsd", "openbsd", "netbsd"
cmd = "xdg-open"
}
*/
