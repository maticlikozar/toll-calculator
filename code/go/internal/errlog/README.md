# errlog

Package that extends the golang error functionality with tracing of where the error has originated.

The intended usage is to wrap errors coming from external packages or newly created errors.
The filename, line number and function name where errlog was called are bundled within the returned error.

When propagating errors between application layers (such as internal packages, repository, service and handler)
calling errlog on each level appends the calling location from each level generating a trace from where the error has originated.

When logging the error with `errlog`, the error trace can be extracted by calling `errlog.StackLog`.
The location is returned as a `err_loc` log field together with any additional custom log fields provided.

The package offers several options for creating, wrapping, merging and comparing errors.

## Usage

The following example demonstrates the error propagation between functions and the logging at the end.

```go
package main

import (
    "toll/internal/log"
    "toll/internal/errlog"
)

var (
    predefinedError = errlog.New("predefiend error")
)

func main() {
    err := f3()

    log.WithFields(
        errlog.StackLog(
            err,
            log.Fields{
                "custom field": "value",
            },
        ),
    ).Errore(err, "error in main")
}

func f1() error {
    return errlog.New("err1")
}

func f2() error {
    return errlog.Error(f1())
}

func f3() error {
    return errlog.Merge(f2(), predefinedError)
}
```

When executed the following log is generated:

```text
{
    "level":"error",
    "error":"predefiend error: err1",
    "custom field":"value",
    "err_loc":[
        "cmd/test/main.go:10 init",
        "cmd/test/main.go:35 f3",
        "cmd/test/main.go:31 f2",
        "cmd/test/main.go:27 f1"
    ],
    "log_loc":"cmd/test/main.go:17 main",
    "time":1680168091865,
    "message":"error in main"
}
```
