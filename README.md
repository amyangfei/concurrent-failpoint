This is a simple demo to show how we can run unit tests in parallel mode with failpoint injection taking effect and without injection race.

The basic premise is once a failpoint is enabled, its injection terms will take effect until we disable the injection explicitly.

## A injection race example

Taking bad_work.go and bad_work_test.go as an example:

bad_work.go defines a function `badWork`, which has two failpoint injections. In bad_work_test.go we have three test cases, `TestBadWorkPath1` and `TestBadWorkPath2` enable a failpoint respectively, parallel mode is enabled in each test case. When we enable failpoint and run test cases, we will find the injection race. Pay attention to the output of `bad work path...`, the injection `cf/bad-path1` will always work after its first injection.

```bash
➜  failpoint-ctl enable .
➜  go test -count=2 -v -run TestBadWork .
=== RUN   TestBadWorkNormal
=== PAUSE TestBadWorkNormal
=== RUN   TestBadWorkPath1
=== PAUSE TestBadWorkPath1
=== RUN   TestBadWorkPath2
=== PAUSE TestBadWorkPath2
=== CONT  TestBadWorkNormal
=== CONT  TestBadWorkPath2
bad work path2
--- PASS: TestBadWorkPath2 (0.00s)
=== CONT  TestBadWorkPath1
bad work path1
--- PASS: TestBadWorkPath1 (0.00s)
bad work path normal
--- PASS: TestBadWorkNormal (0.00s)
=== RUN   TestBadWorkNormal
=== PAUSE TestBadWorkNormal
=== RUN   TestBadWorkPath1
=== PAUSE TestBadWorkPath1
=== RUN   TestBadWorkPath2
=== PAUSE TestBadWorkPath2
=== CONT  TestBadWorkNormal
bad work path1
--- PASS: TestBadWorkNormal (0.00s)
=== CONT  TestBadWorkPath2
=== CONT  TestBadWorkPath1
bad work path1
--- PASS: TestBadWorkPath2 (0.00s)
bad work path1
--- PASS: TestBadWorkPath1 (0.00s)
PASS
ok      cf      0.023s
```

## A workaround for parallel failpoint injection

~In work.go, the failpoint injection path is encoded with a dynamic function name, which is the entrypoint of test case, where we have enabled a failpoint. This time the failpoint injection will only take effect when the matching test case is running.~

We can use InjectContext and failpoint.WithHook to pass a test independent context key value to achieve running parallel test cases without injection race.

```bash
➜  failpoint-ctl enable .
➜  go test -count=2 -v -run TestWork
=== RUN   TestWorkNormal
=== PAUSE TestWorkNormal
=== RUN   TestWorkPath1
=== PAUSE TestWorkPath1
=== RUN   TestWorkPath2
=== PAUSE TestWorkPath2
=== RUN   TestWorkPathPanic
=== PAUSE TestWorkPathPanic
=== CONT  TestWorkNormal
work path normal
--- PASS: TestWorkNormal (0.00s)
=== CONT  TestWorkPathPanic
--- PASS: TestWorkPathPanic (0.00s)
=== CONT  TestWorkPath2
work path2
--- PASS: TestWorkPath2 (0.00s)
=== CONT  TestWorkPath1
work path1
--- PASS: TestWorkPath1 (0.00s)
=== RUN   TestWorkNormal
=== PAUSE TestWorkNormal
=== RUN   TestWorkPath1
=== PAUSE TestWorkPath1
=== RUN   TestWorkPath2
=== PAUSE TestWorkPath2
=== RUN   TestWorkPathPanic
=== PAUSE TestWorkPathPanic
=== CONT  TestWorkNormal
work path normal
--- PASS: TestWorkNormal (0.00s)
=== CONT  TestWorkPathPanic
--- PASS: TestWorkPathPanic (0.00s)
=== CONT  TestWorkPath2
work path2
--- PASS: TestWorkPath2 (0.00s)
=== CONT  TestWorkPath1
work path1
--- PASS: TestWorkPath1 (0.00s)
PASS
ok      cf      0.005s
```

The basic code is as follows

```go
// business code with failpoint injection
func work(ctx context.Context) {
	failpoint.InjectContext(ctx, "path1", func() {
		// do something
	})
	// ...
}

// test case entrypoint, with failpoint enable
func TestWorkPath1(t *testing.T) {
	failpoint.Enable("cf/path1", "return(true)")
	ctx := failpoint.WithHook(context.Background(), func(ctx context.Context, fpname string) bool {
		// only the code is running from this test case contains the context key of `fpname`
		return ctx.Value(fpname) != nil
	})
	ctx = context.WithValue(ctx, "cf/path1", struct{}{})
	work(ctx)
}
```
