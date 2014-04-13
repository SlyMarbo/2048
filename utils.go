// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

// Package terminal provides support functions for dealing with terminals, as
// commonly found on UNIX systems.
//
// Putting a terminal into raw mode is the most common requirement:
//
// 	oldState, err := terminal.MakeRaw(0)
// 	if err != nil {
// 	        panic(err)
// 	}
// 	defer terminal.Restore(0, oldState)
package main

import (
	"fmt"
	"io"
	"syscall"
	"unsafe"
)

type Tile struct {
	Number
	Merged bool
}

type Number uint8

func (n Number) Int() int {
	return 1 << (n + 1)
}

func (n Number) String() string {
	return fmt.Sprint(1 << (n + 1))
}

type Direction uint8

const (
	_              = iota
	Left Direction = iota
	Right
	Up
	Down
)

// x = 0 is left-most, y = 0 is bottom-most
type Index struct {
	X, Y int
}

// Returns success.
func (i *Index) Increment(d Direction, max int) bool {
	switch d {
	case Left:
		if i.X == 0 {
			return false
		}
		i.X--
	case Right:
		if i.X == max-1 {
			return false
		}
		i.X++
	case Up:
		if i.Y == max-1 {
			return false
		}
		i.Y++
	case Down:
		if i.Y == 0 {
			return false
		}
		i.Y--
	}
	return true
}

func (i *Index) Decrement(d Direction, max int) bool {
	switch d {
	case Left:
		if i.X == max {
			return false
		}
		i.X++
	case Right:
		if i.X == 0 {
			return false
		}
		i.X--
	case Up:
		if i.Y == 0 {
			return false
		}
		i.Y--
	case Down:
		if i.Y == max {
			return false
		}
		i.Y++
	}
	return true
}

func (i *Index) Next(d Direction, max int) *Index {
	out := &Index{i.X, i.Y}
	if out.Increment(d, max) {
		return out
	}
	return nil
}

// State contains the state of a terminal.
type State struct {
	termios syscall.Termios
}

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd int) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}

// MakeRaw put the terminal connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be
// restored.
func MakeRaw(fd int) (*State, error) {
	var oldState State
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&oldState.termios)), 0, 0, 0); err != 0 {
		return nil, err
	}

	newState := oldState.termios
	newState.Iflag &^= syscall.ISTRIP | syscall.INLCR | syscall.ICRNL | syscall.IGNCR | syscall.IXON | syscall.IXOFF
	newState.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&newState)), 0, 0, 0); err != 0 {
		return nil, err
	}

	return &oldState, nil
}

// Restore restores the terminal connected to the given file descriptor to a
// previous state.
func Restore(fd int, state *State) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&state.termios)), 0, 0, 0)
	return err
}

// GetSize returns the dimensions of the given terminal.
func GetSize(fd int) (width, height int, err error) {
	var dimensions [4]uint16

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&dimensions)), 0, 0, 0); err != 0 {
		return -1, -1, err
	}
	return int(dimensions[1]), int(dimensions[0]), nil
}

type ReadWriter struct {
	io.Reader
	io.Writer
}
