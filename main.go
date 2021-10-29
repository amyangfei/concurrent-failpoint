package cf

import (
	"fmt"
	"runtime"

	"github.com/pingcap/failpoint"
)

// Return the first function name after removing `skip` callers.
func caller(skip int) string {
	pc := make([]uintptr, 1024)
	n := runtime.Callers(0, pc)
	return runtime.FuncForPC(pc[n-skip-1]).Name()
}

// caller2 alias to caller(2)
// In most cases skip=2, which are
//   - testing.(*T).Run
//   - testing.tRunner()
func caller2() string {
	return caller(2)
}

func work() {
	failpoint.Inject("path1."+caller2(), func() {
		fmt.Println("work path1")
		failpoint.Return()
	})
	failpoint.Inject("path2."+caller2(), func() {
		fmt.Println("work path2")
		failpoint.Return()
	})
	failpoint.Inject("path3."+caller2(), nil)
	fmt.Println("normal worker path")
}
