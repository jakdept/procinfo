package procinfo

import (
	"bufio"
	"bytes"
	"os"

	"github.com/pkg/errors"
)

var (
	ProcLock = "/proc/locks"
)

var (
	ErrPermissionDenied = errors.New("cannot read /proc/locks")
)

type Lock struct {
	Priority  uint64 `json:"lockPriority"`
	Exclusive bool   `json:"exclusive"`
	ByteRange bool   `json:"isByteRange"`
	PID       uint32 `json:"pid"`
	DevMajor  uint16 `json:"fsMajor"`
	DevMinor  uint16 `json:"fsMinor"`
	Inode     uint64 `json:"inode"`
}

func GetAllLocks() ([]Lock, error) {
	f, err := os.Open(ProcLock)
	if err != nil && !os.IsNotExist(err) {
		return []Lock{}, ErrPermissionDenied
	} else if err != nil {
		return []Lock{}, errors.Wrapf(err, "could not read %s", ProcLock)
	}
	defer f.Close()

	var locks []Lock
	var fields [][]byte
	s := bufio.NewScanner(f)
	for s.Scan() {
		fields = bytes.Fields(s.Bytes())

	}
	err = s.Err()

	return locks, nil
}
