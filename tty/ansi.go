package tty

import "strings"

const esc = "\033"

type EraseCommand struct {
	EntireScreen string
}

type ScreenCommand struct {
	EnableAlternativeBuffer  string
	DisableAlternativeBuffer string
}

type AnsiCommand struct {
	ERASE  EraseCommand
	Screen ScreenCommand
}

var erase = EraseCommand{
	EntireScreen: concatCommands(esc, "[2J"),
}

var screen = ScreenCommand{
	EnableAlternativeBuffer:  concatCommands(esc, "[?1049h"),
	DisableAlternativeBuffer: concatCommands(esc, "[?1049l"),
}

var ANSI = &AnsiCommand{
	ERASE:  erase,
	Screen: screen,
}

func concatCommands(commands ...string) string {
	var sb strings.Builder
	for _, cmd := range commands {
		sb.WriteString(cmd)
	}
	return sb.String()
}
