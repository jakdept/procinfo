package procinfo

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

const (
	ProcLock = "/proc/locks"
)

var (
	ErrPermissionDenied = errors.New("cannot read /proc/locks")
)

type Lock struct {
	Priority  uint32 `json:"lockPriority"`
	ByteRange bool   `json:"isByteRange"`
	Exclusive bool   `json:"exclusive"`
	PID       uint32 `json:"pid"`
	DevMajor  uint16 `json:"fsMajor"`
	DevMinor  uint16 `json:"fsMinor"`
	Inode     uint64 `json:"inode"`
}

const (
	ProcLockColPriority = iota
	ProcLockColType
	_
	ProcLockColExclusive
	ProcLockColPID
	ProcLockColInode
	ProcLockColByteStart
	ProcLockColByteEnd
)

var FileLocks = fileLocksType{
	l:           sync.RWMutex{},
	clearSignal: make(chan struct{}),
	dedupWait:   make(chan struct{}),
	delay:       time.Millisecond * 200,
}

type fileLocksType struct {
	locks       *[]Lock
	l           sync.RWMutex
	t           *time.Timer
	delay       time.Duration
	clearSignal chan struct{}
	dedupWait   chan struct{}
}

func (l *fileLocksType) waitForSignal() {
	go func() {
		// multiple instances of this func collapse into one
		// if two goroutines end up in this func, the first two cases will link up
		select {
		case l.dedupWait <- struct{}{}:
			// do nothing here, just kill this thread
			return
		case <-l.dedupWait:
			// relaunch this thread
			l.waitForSignal()
			return
		case <-l.t.C:
			// the timer fired
		case <-l.clearSignal:
			l.t.Stop()
			// you got the signal to cancel, so do it
		}
		l.l.Lock()
		defer l.l.Unlock()
		l.locks = nil
	}()
}

// external clear function - see waitForSignal() for the real one
func (l *fileLocksType) Clear() {
	// use a select with default so it falls through if needed
	select {
	case l.clearSignal <- struct{}{}:
	default:
	}
}

func (l *fileLocksType) isSet() bool {
	l.l.RLock()
	defer l.l.RUnlock()
	return l.locks != nil
}

func (l *fileLocksType) SetExpirationAndClear(t time.Duration) {
	l.delay = t
	l.Clear()
}

func (l *fileLocksType) GetAll() ([]Lock, error) {
	if err := l.load(); err != nil {
		return []Lock{}, err
	}

	l.l.RLock()
	defer l.l.RUnlock()
	return *l.locks, nil
}

func (l *fileLocksType) CheckInode(inodeNum uint64) ([]Process, error) {
	if err := l.load(); err != nil {
		return []Process{}, err
	}

	l.l.RLock()
	defer l.l.RUnlock()

	var processes []Process
	for _, eachLock := range *l.locks {
		if eachLock.Inode == inodeNum {
			processes = append(processes, getProcessByPid(eachLock.PID))
		}
	}
	return processes, nil
}

func (l *fileLocksType) CheckFileInfo(fileinfo os.FileInfo) ([]Process, error) {
	// #todo# check and make sure this will always assert?
	return l.CheckInode(fileinfo.Sys().(syscall.Stat_t).Ino)
}

func (l *fileLocksType) CheckFilePath(path string) ([]Process, error) {
	// #todo# check and make sure this will always assert?
	if finfo, err := os.Stat(path); err != nil {
		return []Process{}, err
	} else {
		return l.CheckFileInfo(finfo)
	}
}

func (l *fileLocksType) load() error {
	if l.isSet() {
		return nil
	}

	l.l.Lock()
	defer l.l.Unlock()

	locks, err := populate()
	if err != nil {
		return err
	}

	// don't actually set these until you know you're setting it
	defer l.waitForSignal()
	if l.t == nil {
		l.t = time.NewTimer(l.delay)
	}
	defer l.t.Reset(l.delay)

	l.locks = &locks
	return nil
}

const FmtProcLocks = "%d: %s %s %s %d %x:%x:%d"

func populate() ([]Lock, error) {
	// open /proc/locks for reading
	f, err := os.Open(testPrefix + ProcLock)
	if err != nil && !os.IsNotExist(err) {
		return []Lock{}, ErrPermissionDenied
	} else if err != nil {
		return []Lock{}, errors.Wrapf(err, "could not read %s", ProcLock)
	}
	defer f.Close()

	// set up a few vars to be reused
	var locks []Lock
	var lockType string
	var lockExclusive string
	s := bufio.NewScanner(f)

	// it would be nice  to do something other than a scanner, but idk protobuf
	// and fmt.Scanf requires a specific spacing - not present.
	for s.Scan() {
		var newLock Lock
		fmt.Sscanf(s.Text(), FmtProcLocks,
			&newLock.Priority,
			&lockType,
			new(string),
			&lockExclusive,
			&newLock.PID,
			&newLock.DevMajor,
			&newLock.DevMinor,
			&newLock.Inode,
		)

		if lockType == "POSIX" && !strings.HasSuffix(s.Text(), "EOF") {
			newLock.ByteRange = true
		}

		if lockExclusive == "WRITE" {
			newLock.Exclusive = true
		}
		locks = append(locks, newLock)
	}

	if err = s.Err(); err != nil {
		return []Lock{}, err
	}
	return locks, nil
}
