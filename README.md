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

In work.go, the failpoint injection path is encoded with a dynamic function name, which is the entrypoint of test case, where we have enabled a failpoint. This time the failpoint injection will only take effect when the matching test case is running.

```bash
➜  failpoint-ctl enable .
➜  go test -count=2 -v -run TestWork
=== RUN   TestWorkPath1
=== PAUSE TestWorkPath1
=== RUN   TestWorkPath2
=== PAUSE TestWorkPath2
=== RUN   TestWorkPathPanic
=== PAUSE TestWorkPathPanic
=== CONT  TestWorkPath1
work path1
--- PASS: TestWorkPath1 (0.00s)
=== CONT  TestWorkPathPanic
--- PASS: TestWorkPathPanic (0.00s)
=== CONT  TestWorkPath2
work path2
--- PASS: TestWorkPath2 (0.00s)
=== RUN   TestWorkPath1
=== PAUSE TestWorkPath1
=== RUN   TestWorkPath2
=== PAUSE TestWorkPath2
=== RUN   TestWorkPathPanic
=== PAUSE TestWorkPathPanic
=== CONT  TestWorkPath1
work path1
--- PASS: TestWorkPath1 (0.00s)
=== CONT  TestWorkPathPanic
--- PASS: TestWorkPathPanic (0.00s)
=== CONT  TestWorkPath2
work path2
--- PASS: TestWorkPath2 (0.00s)
PASS
ok      cf      0.005s
```

The trick code is as follows

```go
// business code with failpoint injection
func() {
	// "path1." is a common injection string
	// caller2() is a dynamic function name, which is often the function name
	// of test case.
	failpoint.Inject("path1."+caller2(), func() {
		// do something
	})
	// ...
}

// test case entrypoint, with failpoint enable
func TestWorkPath1(t *testing.T) {
	// This is still a normal failpoint enable statement, while the last part
	// contains the function name of this test case.
	failpoint.Enable("cf/path1.cf.TestWorkPath1", "return(true)")
	// ...
}
```
