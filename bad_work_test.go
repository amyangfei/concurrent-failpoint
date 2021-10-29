package cf

import (
	"testing"

	"github.com/pingcap/failpoint"
)

func TestBadWorkNormal(t *testing.T) {
	t.Parallel()
	badWork()
}

func TestBadWorkPath1(t *testing.T) {
	t.Parallel()
	failpoint.Enable("cf/bad-path1", "return(true)")
	badWork()
}

func TestBadWorkPath2(t *testing.T) {
	t.Parallel()
	failpoint.Enable("cf/bad-path2", "return(true)")
	badWork()
}
