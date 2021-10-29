package cf

import (
	"testing"

	"github.com/pingcap/failpoint"
	"github.com/stretchr/testify/require"
)

func TestWorkNormal(t *testing.T) {
	t.Parallel()
	work()
}

func TestWorkPath1(t *testing.T) {
	t.Parallel()
	failpoint.Enable("cf/path1.cf.TestWorkPath1", "return(true)")
	work()
}

func TestWorkPath2(t *testing.T) {
	t.Parallel()
	failpoint.Enable("cf/path2.cf.TestWorkPath2", "return(true)")
	work()
}

func TestWorkPathPanic(t *testing.T) {
	t.Parallel()
	failpoint.Enable("cf/path3.cf.TestWorkPathPanic", "panic(`xxx`)")
	require.PanicsWithValue(t, "failpoint panic: xxx", func() { work() })
}
