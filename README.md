# go-test-example-diff

Reads `go test` output, looks for failed examples, and prints a diff for each
one.

## Usage

Given an failing example:

```go
func ExampleMain() {
	fmt.Println(`{
  "name": "bob",
  "date_of_birth": "1970",
}`)
	// Output: {
	//   "name": "bob",
	//   "date_of_birth": "1970"
	// }
}
```

When we execute the tests, it fails with the following output:

```
$ go test ./...
```

```diff
--- FAIL: ExampleMain (0.00s)
got:
{
  "name": "bob",
  "date_of_birth": "1970",
}
want:
{
  "name": "bob",
  "date_of_birth": "1970"
}
FAIL
FAIL    github.com/benjaminheng/go-test-example-diff    0.005s
FAIL
```

Pipe the test output to the program and a diff is added to the output:

```
$ go test ./... | go-test-example-diff
```

```diff
--- FAIL: ExampleMain (0.00s)
got:
{
  "name": "bob",
  "date_of_birth": "1970",
}
want:
{
  "name": "bob",
  "date_of_birth": "1970"
}
diff:
@@ -1,4 +1,4 @@
 {
   "name": "bob",
-  "date_of_birth": "1970",
+  "date_of_birth": "1970"
 }
FAIL
FAIL    github.com/benjaminheng/go-test-example-diff    0.005s
FAIL
```
