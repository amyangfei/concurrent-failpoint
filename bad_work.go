package cf

import (
	"fmt"

	"github.com/pingcap/failpoint"
)

func badWork() {
	failpoint.Inject("bad-path1", func() {
		fmt.Println("bad work path1")
		failpoint.Return()
	})
	failpoint.Inject("bad-path2", func() {
		fmt.Println("bad work path2")
		failpoint.Return()
	})
	fmt.Println("bad work path normal")
}
