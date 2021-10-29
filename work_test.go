package cf

import (
	"context"
	"testing"

	"github.com/pingcap/failpoint"
	"github.com/stretchr/testify/require"
)

func TestWorkNormal(t *testing.T) {
	t.Parallel()
	work(context.Background())
}

func TestWorkPath1(t *testing.T) {
	t.Parallel()
	failpoint.Enable("cf/path1", "return(true)")
	ctx := failpoint.WithHook(context.Background(), func(ctx context.Context, fpname string) bool {
		return ctx.Value(fpname) != nil
	})
	ctx = context.WithValue(ctx, "cf/path1", struct{}{})
	work(ctx)
}

func TestWorkPath2(t *testing.T) {
	t.Parallel()
	failpoint.Enable("cf/path2", "return(true)")
	ctx := failpoint.WithHook(context.Background(), func(ctx context.Context, fpname string) bool {
		return ctx.Value(fpname) != nil
	})
	ctx = context.WithValue(ctx, "cf/path2", struct{}{})
	work(ctx)
}

func TestWorkPathPanic(t *testing.T) {
	t.Parallel()
	failpoint.Enable("cf/path3", "panic(`xxx`)")
	ctx := failpoint.WithHook(context.Background(), func(ctx context.Context, fpname string) bool {
		return ctx.Value(fpname) != nil
	})
	ctx = context.WithValue(ctx, "cf/path3", struct{}{})
	require.PanicsWithValue(t, "failpoint panic: xxx", func() { work(ctx) })
}
