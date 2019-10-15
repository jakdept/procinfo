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

// nolint position           1   2   3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18 19 20 21 22 23 24
const FmtLinuxProcPidStat = "%d (%s) %c %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d" // nolint
var FmtProcStatFile = "/proc/%d/stat"

func (proc *Process) readStat() error {
	f, err := os.Open(fmt.Sprintf(FmtProcStatFile, proc.Pid))
	if err != nil && !os.IsNotExist(err) {
		return ErrPermissionDenied
	} else if err != nil {
		return errors.Wrapf(err, "could not read %s", ProcLock)
	}
	if _, err = fmt.Fscanf(f, FmtLinuxProcPidStat,
		&proc.Pid,
		&proc.OriginalName,
		&proc.State,
		&proc.PPid,
		new(int), // 5
		new(int),
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
		new(int),
		new(int),
		new(int),
		new(int),
		new(int),
		&proc.MemoryUsage, // 24
	); err != nil {
		return err
	}
	return nil
}

var FmtProcCmdlineFile = "/proc/%d/cmdline"

func (proc *Process) readCmdline() error {
	filepath := fmt.Sprintf(FmtProcCmdlineFile, proc.Pid)
	b, err := ioutil.ReadFile(filepath)
	if err != nil && !os.IsNotExist(err) {
		return ErrPermissionDenied
	} else if err != nil {
		return errors.Wrapf(err, "could not read %s", filepath)
	}
	proc.Cmdline = strings.Split(string(b), "\x00")
	return nil
}

var FmtProcEnvFile = "/proc/%d/environ"

func (proc *Process) readEnv() error {
	filepath := fmt.Sprintf(FmtProcEnvFile, proc.Pid)
	b, err := ioutil.ReadFile(filepath)
	if err != nil && !os.IsNotExist(err) {
		return ErrPermissionDenied
	} else if err != nil {
		return errors.Wrapf(err, "could not read %s", filepath)
	}
	proc.Env = strings.Split(string(b), "\x00")
	return nil
}

var FmtProcCwdFile = "/proc/%d/cwd"

func (proc *Process) readCwd() error {
	link, err := os.Readlink(fmt.Sprintf(FmtProcCmdlineFile, proc.Pid))
	if err != nil && !os.IsNotExist(err) {
		return ErrPermissionDenied
	} else if err != nil {
		return errors.Wrapf(err, "could not read %s", ProcLock)
	}
	proc.Cwd = link
	return nil
}

func getProcessByPid(pid uint32) Process {
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
