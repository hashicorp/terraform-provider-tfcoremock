// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

//go:build !windows
// +build !windows

package state_store

import (
	"io"
	"log"
	"os"
	"syscall"
)

// use fcntl POSIX locks for the most consistent behavior across platforms, and
// hopefully some campatibility over NFS and CIFS.
func (f *Fs) lock(state *os.File) error {
	log.Printf("[TRACE] statemgr.Filesystem: locking %s using fcntl flock", state.Name())
	flock := &syscall.Flock_t{
		Type:   syscall.F_RDLCK | syscall.F_WRLCK,
		Whence: int16(io.SeekStart),
		Start:  0,
		Len:    0,
	}

	fd := state.Fd()
	return syscall.FcntlFlock(fd, syscall.F_SETLK, flock)
}

func (f *Fs) unlock(state *os.File) error {
	log.Printf("[TRACE] statemgr.Filesystem: unlocking %s using fcntl flock", state.Name())
	flock := &syscall.Flock_t{
		Type:   syscall.F_UNLCK,
		Whence: int16(io.SeekStart),
		Start:  0,
		Len:    0,
	}

	fd := state.Fd()
	return syscall.FcntlFlock(fd, syscall.F_SETLK, flock)
}
