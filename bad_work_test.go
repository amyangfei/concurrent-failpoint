package cf

import (
	"testing"

	"github.com/pingcap/failpoint"
)

func TestNormalBadWork(t *testing.T) {
	t.Parallel()
	badWork()
}

func TestWorkBadPath1(t *testing.T) {
	t.Parallel()
	failpoint.Enable("cf/bad-path1", "return(true)")
	badWork()
}

func TestWorkBadPath2(t *testing.T) {
	t.Parallel()
	failpoint.Enable("cf/bad-path2", "return(true)")
	badWork()
}
