package errors

import (
    "fmt"
    "path"
    "runtime"
)

func wrapperErr(text string) (output string) {
    _, file, line, ok := runtime.Caller(2)
    if !ok {
        file = "???"
        line = 0
    }
    output = fmt.Sprintf("%s:%d %s", path.Base(file), line, text)
    return
}

func New(text string) (err error) {
    //buf := make([]byte, 1204)
    //length := runtime.Stack(buf, false)
    errImpl := new(dnsError)
    errImpl.err = wrapperErr(text)
    //errImpl.stack = string(buf[:length])
    return errImpl
}

//type DnsError interface {
//Error() string
//Stack() string
//}

type dnsError struct {
    err string
    //stack string
}

func (e *dnsError) Error() string {
    return e.err
}

//func (e *dnsError) Stack() string {
//return e.stack
//}
