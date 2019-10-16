package procinfo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllLocks(t *testing.T) {
	testResult, err := FileLocks.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, Testdata_Lock, testResult)
}

func TestCheckInode(t *testing.T) {
	for inode, expectedProcesses := range TestData_CheckInode_Processes {
		t.Run(fmt.Sprintf("%s %d", t.Name(), inode), func(t *testing.T) {
			testResult, err := FileLocks.CheckInode(inode)
			assert.NoError(t, err)
			assert.Equal(t, expectedProcesses, testResult)
		})
	}
}
