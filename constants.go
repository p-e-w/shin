// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2022  Philipp Emanuel Weidmann <pew@worldwidemann.com>

package main

// https://github.com/ibus/ibus/blob/main/src/ibustypes.h
const (
	ShiftMask    = 1 << 0
	LockMask     = 1 << 1
	ControlMask  = 1 << 2
	Mod1Mask     = 1 << 3
	Mod2Mask     = 1 << 4
	Mod3Mask     = 1 << 5
	Mod4Mask     = 1 << 6
	Mod5Mask     = 1 << 7
	Button1Mask  = 1 << 8
	Button2Mask  = 1 << 9
	Button3Mask  = 1 << 10
	Button4Mask  = 1 << 11
	Button5Mask  = 1 << 12
	HandledMask  = 1 << 24
	ForwardMask  = 1 << 25
	IgnoredMask  = ForwardMask
	SuperMask    = 1 << 26
	HyperMask    = 1 << 27
	MetaMask     = 1 << 28
	ReleaseMask  = 1 << 30
	ModifierMask = 0x5f001fff
)

// https://github.com/ibus/ibus/blob/main/src/ibuskeysyms.h
const (
	KeyBackSpace  = 0xff08
	KeyTab        = 0xff09
	KeyReturn     = 0xff0d
	KeyEscape     = 0xff1b
	KeyDelete     = 0xffff
	KeyHome       = 0xff50
	KeyLeft       = 0xff51
	KeyUp         = 0xff52
	KeyRight      = 0xff53
	KeyDown       = 0xff54
	KeyPageUp     = 0xff55
	KeyPageDown   = 0xff56
	KeyEnd        = 0xff57
	KeyInsert     = 0xff63
	KeyKPTab      = 0xff89
	KeyKPEnter    = 0xff8d
	KeyKPHome     = 0xff95
	KeyKPLeft     = 0xff96
	KeyKPUp       = 0xff97
	KeyKPRight    = 0xff98
	KeyKPDown     = 0xff99
	KeyKPPageUp   = 0xff9a
	KeyKPPageDown = 0xff9b
	KeyKPEnd      = 0xff9c
	KeyKPInsert   = 0xff9e
	KeyKPDelete   = 0xff9f
)

const (
	engineBasePath = "/org/freedesktop/IBus/Engine/Shin"
	busName        = "org.freedesktop.IBus.Shin"
)
