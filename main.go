// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2022  Philipp Emanuel Weidmann <pew@worldwidemann.com>

package main

import (
	"fmt"

	"github.com/BambooEngine/goibus/ibus"
	"github.com/godbus/dbus"
)

const engineBasePath = "/org/freedesktop/IBus/Engine/Shin"
const busName = "org.freedesktop.IBus.Shin"

type engine struct {
	ibus.Engine
}

func (e *engine) ProcessKeyEvent(keyval uint32, keycode uint32, state uint32) (bool, *dbus.Error) {
	fmt.Printf("Shin: ProcessKeyEvent(keyval = %v, keycode = %v, state = %v)\n", keyval, keycode, state)

	return false, nil
}

func main() {
	fmt.Println("Shin: Starting")

	bus := ibus.NewBus()
	connection := bus.GetDbusConn()

	engineId := 0

	ibus.NewFactory(connection, func(connection *dbus.Conn, engineName string) dbus.ObjectPath {
		engineId++

		path := dbus.ObjectPath(fmt.Sprintf("%v/%v", engineBasePath, engineId))
		engine := &engine{ibus.BaseEngine(connection, path)}

		ibus.PublishEngine(connection, path, engine)

		return path
	})

	bus.RequestName(busName, 0)

	fmt.Println("Shin: Started")

	select {}
}
