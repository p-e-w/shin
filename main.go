// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2022  Philipp Emanuel Weidmann <pew@worldwidemann.com>

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
)

// TODO: This only matches ASCII word boundaries!
var wordBoundaryRegex = regexp.MustCompile(`\b`)

func bytePosToCharacterPos(text string, bytePos int) int {
	if bytePos == len(text) {
		return utf8.RuneCountInString(text)
	}

	characterPos := 0

	for i := range text {
		if i == bytePos {
			return characterPos
		}

		characterPos++
	}

	panic("byte position does not correspond to a character position")
}

type engine struct {
	ibus.Engine

	text      string
	cursorPos uint32

	history       *history
	historyPrefix string
	historyIndex  uint32
}

func (e *engine) exit() {
	log.Printf("Exiting")

	go func() {
		// Give the calling function time to return,
		// so that pending IBus operations can finish.
		time.Sleep(100 * time.Millisecond)

		os.Exit(0)
	}()
}

func (e *engine) textLength() uint32 {
	return uint32(utf8.RuneCountInString(e.text))
}

func (e *engine) updateText() {
	log.Printf("updateText(text = '%v', cursorPos = %v)", e.text, e.cursorPos)

	ibusText := ibus.NewText(e.text)
	ibusText.AppendAttr(ibus.IBUS_ATTR_TYPE_UNDERLINE, ibus.IBUS_ATTR_UNDERLINE_SINGLE, 0, e.textLength())

	e.UpdatePreeditText(ibusText, e.cursorPos, e.text != "")
}

func (e *engine) clearText() {
	e.text = ""
	e.cursorPos = 0
	e.updateText()
	e.historyIndex = 0
}

func (e *engine) previousWordBoundary() uint32 {
	text := string([]rune(e.text)[:e.cursorPos])

	locations := wordBoundaryRegex.FindAllStringIndex(text, -1)

	if locations == nil {
		return 0
	}

	boundary := locations[len(locations)-1][1]

	if boundary == len(text) {
		if len(locations) > 1 {
			boundary = locations[len(locations)-2][1]
		} else {
			return 0
		}
	}

	return uint32(bytePosToCharacterPos(text, boundary))
}

func (e *engine) nextWordBoundary() uint32 {
	text := string([]rune(e.text)[e.cursorPos:])

	locations := wordBoundaryRegex.FindAllStringIndex(text, 2)

	if locations == nil {
		return e.textLength()
	}

	boundary := locations[0][0]

	if boundary == 0 {
		if len(locations) > 1 {
			boundary = locations[1][0]
		} else {
			return e.textLength()
		}
	}

	return e.cursorPos + uint32(bytePosToCharacterPos(text, boundary))
}

func (e *engine) moveCursor(offset int32) {
	newCursorPos := int32(e.cursorPos) + offset

	textLength := e.textLength()

	if newCursorPos < 0 {
		e.cursorPos = 0
	} else if newCursorPos > int32(textLength) {
		e.cursorPos = textLength
	} else {
		e.cursorPos = uint32(newCursorPos)
	}
}

func (e *engine) ProcessKeyEvent(keyval uint32, keycode uint32, state uint32) (bool, *dbus.Error) {
	log.Printf("ProcessKeyEvent(keyval = %v, keycode = %v, state = %v)", keyval, keycode, state)

	if state&ReleaseMask != 0 {
		// Key released.
		return true, nil
	}

	switch keyval {
	case KeyReturn, KeyKPEnter:
		if e.text == "" {
			e.exit()
			return true, nil
		}

		err := e.history.addCommand(e.text)

		if err != nil {
			log.Printf("Error from addCommand: %v", err)
		}

		command := exec.Command("bash", "-c", e.text)
		output, err := command.CombinedOutput()

		_, isExitError := err.(*exec.ExitError)

		e.clearText()

		if err == nil || isExitError {
			// Remove leading and trailing newlines from output
			// to allow single-line output to flow into the
			// surrounding text ...
			outputText := strings.Trim(string(output), "\n")

			// ... but place multi-line output in a separate
			// block surrounded by newlines, so that tabular
			// alignment is preserved regardless of where
			// the output is inserted.
			if strings.Contains(outputText, "\n") {
				outputText = "\n" + outputText + "\n"
			}

			e.CommitText(ibus.NewText(outputText))
		} else {
			log.Printf("Error from CombinedOutput: %v", err)
		}

		e.exit()

		return true, nil

	case KeyEscape:
		if e.historyIndex == 0 || e.historyPrefix == "" {
			// Exit input engine.
			e.clearText()
			e.exit()
		} else {
			// Exit history mode.
			e.text = e.historyPrefix
			e.cursorPos = e.textLength()
			e.updateText()
			e.historyIndex = 0
		}
		return true, nil

	case KeyBackSpace:
		if e.cursorPos > 0 {
			characters := []rune(e.text)
			characters = append(characters[:e.cursorPos-1], characters[e.cursorPos:]...)
			e.text = string(characters)

			e.cursorPos--

			e.updateText()
		}
		e.historyIndex = 0
		return true, nil

	case KeyDelete, KeyKPDelete:
		if e.cursorPos < e.textLength() {
			characters := []rune(e.text)
			characters = append(characters[:e.cursorPos], characters[e.cursorPos+1:]...)
			e.text = string(characters)

			e.updateText()
		}
		e.historyIndex = 0
		return true, nil

	case KeyUp, KeyKPUp:
		if e.historyIndex == 0 {
			e.historyPrefix = e.text
		}

		command, err := e.history.getRecentCommand(e.historyPrefix, e.historyIndex)

		if err == nil {
			e.text = command
			e.cursorPos = e.textLength()
			e.updateText()
			e.historyIndex++
		} else if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Error from getRecentCommand: %v", err)
		}

		return true, nil

	case KeyDown, KeyKPDown:
		if e.historyIndex == 0 {
			return true, nil
		}

		e.historyIndex--

		command := e.historyPrefix

		if e.historyIndex > 0 {
			var err error
			command, err = e.history.getRecentCommand(e.historyPrefix, e.historyIndex-1)

			if err != nil {
				log.Printf("Error from getRecentCommand: %v", err)
				e.historyIndex++
				return true, nil
			}
		}

		e.text = command
		e.cursorPos = e.textLength()
		e.updateText()

		return true, nil

	case KeyLeft, KeyKPLeft:
		if state&ControlMask != 0 {
			// Ctrl+Left.
			e.cursorPos = e.previousWordBoundary()
		} else {
			e.moveCursor(-1)
		}
		e.updateText()
		e.historyIndex = 0
		return true, nil

	case KeyRight, KeyKPRight:
		if state&ControlMask != 0 {
			// Ctrl+Right.
			e.cursorPos = e.nextWordBoundary()
		} else {
			e.moveCursor(1)
		}
		e.updateText()
		e.historyIndex = 0
		return true, nil

	case KeyHome, KeyKPHome:
		e.cursorPos = 0
		e.updateText()
		e.historyIndex = 0
		return true, nil

	case KeyEnd, KeyKPEnd:
		e.cursorPos = e.textLength()
		e.updateText()
		e.historyIndex = 0
		return true, nil
	}

	character := rune(keyval)

	if character <= unicode.MaxLatin1 && unicode.IsPrint(character) {
		if e.cursorPos == e.textLength() {
			e.text += string(character)
		} else {
			characters := []rune(e.text)
			characters = append(characters[:e.cursorPos+1], characters[e.cursorPos:]...)
			characters[e.cursorPos] = character
			e.text = string(characters)
		}

		e.cursorPos++
		e.updateText()
		e.historyIndex = 0

		return true, nil
	}

	return true, nil
}

func (e *engine) FocusOut() *dbus.Error {
	e.clearText()
	e.exit()

	return nil
}

func main() {
	log.SetPrefix("[Shin] ")

	log.Printf("Starting")

	history, err := openHistory()

	if err != nil {
		log.Fatalf("Error from openHistory: %v", err)
	}

	bus := ibus.NewBus()
	connection := bus.GetDbusConn()

	engineId := 0

	ibus.NewFactory(connection, func(connection *dbus.Conn, engineName string) dbus.ObjectPath {
		engineId++

		path := dbus.ObjectPath(fmt.Sprintf("%v/%v", engineBasePath, engineId))

		engine := &engine{
			Engine:  ibus.BaseEngine(connection, path),
			history: history,
		}

		ibus.PublishEngine(connection, path, engine)

		return path
	})

	_, err = bus.RequestName(busName, 0)

	if err != nil {
		log.Fatalf("Error from RequestName: %v", err)
	}

	log.Printf("Started")

	select {}
}
