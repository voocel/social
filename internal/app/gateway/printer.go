package gateway

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/mattn/go-runewidth"
	"social/pkg/version"
)

type Colors struct {
	Black   string
	Red     string
	Green   string
	Yellow  string
	Blue    string
	Magenta string
	Cyan    string
	White   string
	Reset   string
}

var DefaultColors = Colors{
	Black:   "\u001b[90m",
	Red:     "\u001b[91m",
	Green:   "\u001b[92m",
	Yellow:  "\u001b[93m",
	Blue:    "\u001b[94m",
	Magenta: "\u001b[95m",
	Cyan:    "\u001b[96m",
	White:   "\u001b[97m",
	Reset:   "\u001b[0m",
}

// `Hello, world \ʕ◔ϖ◔ʔ/`
func startupMessage(addr string, name string) {
	colors := DefaultColors
	value := func(s string, width int) string {
		pad := width - len(s)
		str := ""
		for i := 0; i < pad; i++ {
			str += "."
		}
		if s == "Disabled" {
			str += " " + s
		} else {
			str += fmt.Sprintf(" %s%s%s", colors.Cyan, s, colors.Black)
		}
		return str
	}

	center := func(s string, width int) string {
		pad := strconv.Itoa((width - len(s)) / 2)
		str := fmt.Sprintf("%"+pad+"s", " ")
		str += s
		str += fmt.Sprintf("%"+pad+"s", " ")
		if len(str) < width {
			str += " "
		}
		return str
	}

	centerValue := func(s string, width int) string {
		pad := strconv.Itoa((width - runewidth.StringWidth(s)) / 2)
		str := fmt.Sprintf("%"+pad+"s", " ")
		str += fmt.Sprintf("%s%s%s", colors.Cyan, s, colors.Black)
		str += fmt.Sprintf("%"+pad+"s", " ")
		if runewidth.StringWidth(s)-10 < width && runewidth.StringWidth(s)%2 == 0 {
			// 如果为偶数则在后面添加空格
			str += " "
		}
		return str
	}

	host, port := parseAddr(addr)
	if host == "" {
		host = "0.0.0.0"
	}

	scheme := "http"
	isPrefork := "Disabled"
	proc := strconv.Itoa(runtime.GOMAXPROCS(0))

	mainLogo := colors.Black + " ┌───────────────────────────────────────────────────┐\n"
	if name != "" {
		mainLogo += " │ " + centerValue(name, 49) + " │\n"
	}
	mainLogo += " │ " + centerValue("Social v"+version.VERSION, 49) + " │\n"

	if host == "0.0.0.0" {
		mainLogo +=
			" │ " + center(fmt.Sprintf("%s://127.0.0.1:%s", scheme, port), 49) + " │\n" +
				" │ " + center(fmt.Sprintf("(bound on host 0.0.0.0 and port %s)", port), 49) + " │\n"
	} else {
		mainLogo +=
			" │ " + center(fmt.Sprintf("%s://%s:%s", scheme, host, port), 49) + " │\n"
	}

	mainLogo += fmt.Sprintf(
		" │                                                   │\n"+
			" │ Handlers %s  Processes %s │\n"+
			" │ Prefork .%s  PID ....%s │\n"+
			" └───────────────────────────────────────────────────┘"+
			colors.Reset,
		value(strconv.Itoa(666), 14), value(proc, 12),
		value(isPrefork, 14), value(strconv.Itoa(os.Getpid()), 14),
	)

	var childPidLogo string
	splitMainLogo := strings.Split(mainLogo, "\n")
	splitChildPidLogo := strings.Split(childPidLogo, "\n")

	mainLen := len(splitMainLogo)
	childLen := len(splitChildPidLogo)

	if mainLen > childLen {
		diff := mainLen - childLen
		for i := 0; i < diff; i++ {
			splitChildPidLogo = append(splitChildPidLogo, "")
		}
	} else {
		diff := childLen - mainLen
		for i := 0; i < diff; i++ {
			splitMainLogo = append(splitMainLogo, "")
		}
	}

	output := "\n"
	for i := range splitMainLogo {
		output += colors.Black + splitMainLogo[i] + " " + splitChildPidLogo[i] + "\n"
	}

	out := colorable.NewColorableStdout()
	if os.Getenv("TERM") == "dumb" || os.Getenv("NO_COLOR") == "1" || (!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd())) {
		out = colorable.NewNonColorable(os.Stdout)
	}

	_, _ = fmt.Fprintln(out, output)
}

func parseAddr(raw string) (host, port string) {
	if i := strings.LastIndex(raw, ":"); i != -1 {
		return raw[:i], raw[i+1:]
	}
	return raw, ""
}
