package procinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllLocks(t *testing.T) {
	testResult, err := FileLocks.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, Testdata_Lock, testResult)
}
