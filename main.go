// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2022  Philipp Emanuel Weidmann <pew@worldwidemann.com>

package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
)

type engine struct {
	ibus.Engine
	text      string
	cursorPos uint32
}

func (e *engine) exit() {
	fmt.Println("Shin: Exiting")

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
	fmt.Printf("Shin: updateText(text = '%v', cursorPos = %v)\n", e.text, e.cursorPos)

	ibusText := ibus.NewText(e.text)
	ibusText.AppendAttr(ibus.IBUS_ATTR_TYPE_UNDERLINE, ibus.IBUS_ATTR_UNDERLINE_SINGLE, 0, e.textLength())

	e.UpdatePreeditText(ibusText, e.cursorPos, e.text != "")
}

func (e *engine) clearText() {
	e.text = ""
	e.cursorPos = 0
	e.updateText()
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
	fmt.Printf("Shin: ProcessKeyEvent(keyval = %v, keycode = %v, state = %v)\n", keyval, keycode, state)

	if state&ReleaseMask != 0 {
		// Key released.
		return false, nil
	}

	switch keyval {
	case KeyReturn, KeyKPEnter:
		command := exec.Command("bash", "-c", e.text)
		output, err := command.CombinedOutput()

		_, isExitError := err.(*exec.ExitError)

		e.clearText()

		if err == nil || isExitError {
			e.CommitText(ibus.NewText(string(output)))
		} else {
			fmt.Printf("Shin: Error: %v\n", err)
		}

		e.exit()

		return true, nil

	case KeyEscape:
		e.clearText()
		e.exit()
		return true, nil

	case KeyBackSpace:
		if e.cursorPos > 0 {
			characters := []rune(e.text)
			characters = append(characters[:e.cursorPos-1], characters[e.cursorPos:]...)
			e.text = string(characters)

			e.cursorPos--

			e.updateText()
		}
		return true, nil

	case KeyDelete, KeyKPDelete:
		if e.cursorPos < e.textLength() {
			characters := []rune(e.text)
			characters = append(characters[:e.cursorPos], characters[e.cursorPos+1:]...)
			e.text = string(characters)

			e.updateText()
		}
		return true, nil

	case KeyLeft, KeyKPLeft:
		e.moveCursor(-1)
		e.updateText()
		return true, nil

	case KeyRight, KeyKPRight:
		e.moveCursor(1)
		e.updateText()
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

		return true, nil
	}

	return false, nil
}

func (e *engine) FocusOut() *dbus.Error {
	e.clearText()
	e.exit()

	return nil
}

func main() {
	fmt.Println("Shin: Starting")

	bus := ibus.NewBus()
	connection := bus.GetDbusConn()

	engineId := 0

	ibus.NewFactory(connection, func(connection *dbus.Conn, engineName string) dbus.ObjectPath {
		engineId++

		path := dbus.ObjectPath(fmt.Sprintf("%v/%v", engineBasePath, engineId))
		engine := &engine{ibus.BaseEngine(connection, path), "", 0}

		ibus.PublishEngine(connection, path, engine)

		return path
	})

	bus.RequestName(busName, 0)

	fmt.Println("Shin: Started")

	select {}
}
