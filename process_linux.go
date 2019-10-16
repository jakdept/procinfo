package procinfo

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Process struct {
	// process number and parent process number
	Pid          uint32
	PPid         uint32
	OriginalName string
	State        rune

	UserTime   uint64
	KernelTime uint64

	Nice        uint8
	MemoryUsage uint64

	Cmdline []string
	Env     []string
	Cwd     string
}

const (
	ProcStateRunning      = 'R'
	ProcStateSleeping     = 'S'
	ProcStateWaiting      = 'D'
	ProcStateZombie       = 'Z'
	ProcStateStopped      = 'T'
	ProcStateTraceStopped = 't'
	ProcStatePaging       = 'W'
	ProcStateLegacyDead   = 'X'
	ProcStateDead         = 'x'
	ProcStateWakekill     = 'K'
	ProcStateWaking       = 'W'
	ProcStateParked       = 'P'
)

var testPrefix = ""

// position             1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18
const FmtProcPidStat = "%d %s %c %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d"

// func init() {
// 	// position        1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17
// 	FmtProcPidStat += "%d %s %c %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d"
// 	// position        18 19 20 21 22 23 24 25
// 	FmtProcPidStat += "%d %d %d %d %d %d %d %d"
// 	// // position        18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34
// 	// FmtProcPidStat += "%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d "
// 	// // position        35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52
// 	// FmtProcPidStat += "%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d"
// }

const FmtProcStatFile = "/proc/%d/stat"

func (proc *Process) readStat() error {
	f, err := os.Open(fmt.Sprintf(testPrefix+FmtProcStatFile, proc.Pid))
	if err != nil && !os.IsNotExist(err) {
		return ErrPermissionDenied
	} else if err != nil {
		return errors.Wrapf(err, "could not read %s", ProcLock)
	}
	var procname string
	if _, err = fmt.Fscanf(f, FmtProcPidStat,
		&proc.Pid,   // 1
		&procname,   // 2
		&proc.State, // 3
		&proc.PPid,  // 4
		new(int),
		new(int),
		new(int),
		new(int),
		new(int),
		new(int),
		new(int),
		new(int),
		&proc.UserTime, new(int), // 15
		&proc.KernelTime, new(int), // 16
		new(int),
		&proc.Nice, // 18
	); err != nil {
		return err
	}
	proc.OriginalName = strings.TrimSuffix(strings.TrimPrefix(procname, "("), ")")

	return nil
}

const FmtProcCmdlineFile = "/proc/%d/cmdline"

func (proc *Process) readCmdline() error {
	filepath := fmt.Sprintf(testPrefix+FmtProcCmdlineFile, proc.Pid)
	b, err := ioutil.ReadFile(filepath)
	if err != nil && !os.IsNotExist(err) {
		return ErrPermissionDenied
	} else if err != nil {
		return errors.Wrapf(err, "could not read %s", filepath)
	}
	proc.Cmdline = strings.Split(string(b), "\x00")
	return nil
}

const FmtProcEnvFile = "/proc/%d/environ"

func (proc *Process) readEnv() error {
	filepath := fmt.Sprintf(testPrefix+FmtProcEnvFile, proc.Pid)
	b, err := ioutil.ReadFile(filepath)
	if err != nil && !os.IsNotExist(err) {
		return ErrPermissionDenied
	} else if err != nil {
		return errors.Wrapf(err, "could not read %s", filepath)
	}
	proc.Env = strings.Split(string(b), "\x00")
	return nil
}

const FmtProcCwdFile = "/proc/%d/cwd"

func (proc *Process) readCwd() error {
	filepath := fmt.Sprintf(testPrefix+FmtProcCwdFile, proc.Pid)
	link, err := os.Readlink(filepath)
	if err != nil && !os.IsNotExist(err) {
		return ErrPermissionDenied
	} else if err != nil {
		return errors.Wrapf(err, "could not read %s", filepath)
	}
	proc.Cwd = link
	return nil
}

func GetProcessByPid(pid uint32) Process {
	proc := Process{
		Pid: pid,
	}
	var err error
	if err = proc.readStat(); err != nil {
		return Process{}
	}
	if err = proc.readCwd(); err != nil {
		return Process{}
	}
	if err = proc.readEnv(); err != nil {
		return Process{}
	}
	if err = proc.readCmdline(); err != nil {
		return Process{}
	}
	return proc
}
