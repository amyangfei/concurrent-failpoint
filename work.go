package cf

import (
	"context"
	"fmt"

	"github.com/pingcap/failpoint"
)

func work(ctx context.Context) {
	failpoint.InjectContext(ctx, "path1", func() {
		fmt.Println("work path1")
		failpoint.Return()

	})
	failpoint.InjectContext(ctx, "path2", func() {
		fmt.Println("work path2")
		failpoint.Return()
	})
	failpoint.InjectContext(ctx, "path3", nil)
	fmt.Println("work path normal")
}
