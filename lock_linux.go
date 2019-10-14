package procinfo

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
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

func GetAllLocks() ([]Lock, error) {

	// open /proc/locks for reading
	f, err := os.Open(ProcLock)
	if err != nil && !os.IsNotExist(err) {
		return []Lock{}, ErrPermissionDenied
	} else if err != nil {
		return []Lock{}, errors.Wrapf(err, "could not read %s", ProcLock)
	}
	defer f.Close()

	// set up a few vars to be reused
	var locks []Lock
	var fields []string
	var newLock Lock
	var tempInt int
	s := bufio.NewScanner(f)

	// it would be nice  to do something other than a scanner, but idk protobuf
	// and fmt.Scanf requires a specific spacing - not present.
	for s.Scan() {
		fields = strings.Fields(s.Text())

		// determine the priority of the lock (lowest priority wins)
		tempInt, err = strconv.Atoi(strings.TrimRight(fields[ProcLockColPriority], ":"))
		if err != nil {
			continue
		}
		newLock.Priority = uint32(tempInt)

		switch fields[ProcLockColType] {
		case "POSIX": // POSIX is the only byte-range lock type on Linux
			newLock.ByteRange = true
		default:
			newLock.ByteRange = false
		}

		switch fields[ProcLockColExclusive] {
		case "WRITE":
			newLock.Exclusive = true
		default:
			newLock.Exclusive = false
		}

		tempInt, err = strconv.Atoi(fields[ProcLockColPID])
		if err != nil {
			continue
		}
		newLock.PID = uint32(tempInt)

		fields = strings.Split(fields[ProcLockColInode], ":")
		tempInt, err = strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		newLock.DevMajor = uint16(tempInt)

		tempInt, err = strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		newLock.DevMinor = uint16(tempInt)

		tempInt, err = strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		newLock.Inode = uint64(tempInt)
	}

	if err = s.Err(); err != nil {
		return []Lock{}, err
	}
	return locks, nil
}
