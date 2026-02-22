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
	ERASE    EraseCommand
	Screen   ScreenCommand
	Mouse    MouseEventCommand
	Keyboard KeyboardEventCommand
}

type MouseEventCommand struct {
	EnableMouseEvent  string
	DisableMouseEvent string
}

type KeyboardEventCommand struct {
	EnableKeyboardEvent  string
	DisableKeyboardEvent string
}

var erase = EraseCommand{
	EntireScreen: concatCommands(esc, "[2J"),
}

var screen = ScreenCommand{
	EnableAlternativeBuffer:  concatCommands(esc, "[?1049h"),
	DisableAlternativeBuffer: concatCommands(esc, "[?1049l"),
}

var mouse = MouseEventCommand{
	EnableMouseEvent:  concatCommands(esc, "[?1004h"),
	DisableMouseEvent: concatCommands(esc, "[?1004l"),
}

var keyboard = KeyboardEventCommand{
	EnableKeyboardEvent:  concatCommands(esc, "[?1002h"), //, esc, "[?1006h"),
	DisableKeyboardEvent: concatCommands(esc, "[?1002l"), // esc, "[?1006l"),
}

var ANSI = &AnsiCommand{
	ERASE:    erase,
	Screen:   screen,
	Mouse:    mouse,
	Keyboard: keyboard,
}

func concatCommands(commands ...string) string {
	var sb strings.Builder
	for _, cmd := range commands {
		sb.WriteString(cmd)
	}
	return sb.String()
}
