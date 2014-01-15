package main

import (
    "fmt"
    "runtime"
)

func wrapperErr(text, wrapper string) (output string) {
    _, file, line, ok := runtime.Caller(2)
    if !ok {
        file = "???"
        line = 0
    }
    output = fmt.Sprintf("%s %s:%d %s", wrapper, file, line, text)
    return
}

func NewDnsError(text string) (err error) {
    return &DnsError{wrapperErr("dns", text)}
}

type DnsError struct {
    err string
}

func (e *DnsError) Error() string {
    return e.err
}

type UnimplementedError struct {
    DnsError
}
