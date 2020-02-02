package pick

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"unicode"
)

// OpenBrowser open url each platform default browser
func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Fatal(err)
	}
}

type fder interface {
	Fd() uintptr
}

func CaseInsensitiveIndexFunc(r rune) func(rune) bool {
	lr := unicode.ToUpper(r)

	return func(v rune) bool {
		return lr == unicode.ToUpper(v)
	}
}

func ContainsUpper(query string) bool {
	for _, c := range query {
		if unicode.IsUpper(c) {
			return true
		}
	}
	return false
}

// Global var used to strips ansi sequences
var reANSIEscapeChars = regexp.MustCompile("\x1B\\[(?:[0-9]{1,2}(?:;[0-9]{1,2})?)*[a-zA-Z]")

// Function who strips ansi sequences
func StripANSISequence(s string) string {
	return reANSIEscapeChars.ReplaceAllString(s, "")
}

type causer interface {
	Cause() error
}

type ignorable interface {
	Ignorable() bool
}

type collectResults interface {
	CollectResults() bool
}

type exitStatuser interface {
	ExitStatus() int
}

func IsIgnorableError(err error) bool {
	for e := err; e != nil; {
		switch e.(type) {
		case ignorable:
			return e.(ignorable).Ignorable()
		case causer:
			e = e.(causer).Cause()
		default:
			return false
		}
	}
	return false
}

func IsCollectResultsError(err error) bool {
	for e := err; e != nil; {
		switch e.(type) {
		case collectResults:
			return e.(collectResults).CollectResults()
		case causer:
			e = e.(causer).Cause()
		default:
			return false
		}
	}

	return false
}

func GetExitStatus(err error) (int, bool) {
	for e := err; e != nil; {
		if ese, ok := e.(exitStatuser); ok {
			return ese.ExitStatus(), true
		}
		if cerr, ok := e.(causer); ok {
			e = cerr.Cause()
			continue
		}
		break
	}

	return 1, false
}
