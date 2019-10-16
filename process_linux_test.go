package procinfo

import (
	"fmt"
	"os"
	"testing"

	"github.com/sebdah/goldie"
	"github.com/stretchr/testify/assert"
)

func init() {
	goldie.FixtureDir = "testdata/golden"
}

func TestMain(m *testing.M) {
	testPrefix = "testdata/fixtures"
	defer func() {
		testPrefix = ""
	}()

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestReadStat(t *testing.T) {
	for _, expected := range Testdata_Process {
		t.Run(fmt.Sprintf("%s %d", t.Name(), expected.Pid), func(t *testing.T) {
			testProc := Process{
				Pid: expected.Pid,
			}
			err := testProc.readStat()
			expected.Cwd = ""
			expected.Env = nil
			expected.Cmdline = nil

			assert.NoError(t, err)
			assert.Equal(t, expected, testProc)
		})
	}
}

func TestReadCwd(t *testing.T) {
	for _, expected := range Testdata_Process {
		t.Run(fmt.Sprintf("%s %d", t.Name(), expected.Pid), func(t *testing.T) {
			testProc := Process{
				Pid: expected.Pid,
			}
			err := testProc.readCwd()
			assert.NoError(t, err)
			assert.Equal(t, expected.Cwd, testProc.Cwd)
		})
	}
}

func TestReadCmdline(t *testing.T) {
	for _, expected := range Testdata_Process {
		t.Run(fmt.Sprintf("%s %d", t.Name(), expected.Pid), func(t *testing.T) {
			testProc := Process{
				Pid: expected.Pid,
			}
			err := testProc.readCmdline()
			assert.NoError(t, err)
			assert.Equal(t, expected.Cmdline, testProc.Cmdline)
		})
	}
}

func TestReadEnv(t *testing.T) {
	for _, expected := range Testdata_Process {
		t.Run(fmt.Sprintf("%s %d", t.Name(), expected.Pid), func(t *testing.T) {
			testProc := Process{
				Pid: expected.Pid,
			}
			err := testProc.readEnv()
			assert.NoError(t, err)
			assert.Equal(t, expected.Env, testProc.Env)
		})
	}
}

func TestReadPid(t *testing.T) {
	for id, expected := range Testdata_Process {
		t.Run(fmt.Sprintf("%s %d", t.Name(), id), func(t *testing.T) {
			testProc := GetProcessByPid(expected.Pid)
			assert.Equal(t, expected, testProc)
		})
	}
}
