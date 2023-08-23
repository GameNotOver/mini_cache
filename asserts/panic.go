package asserts

import (
	"errors"
	"fmt"
)

func doPanic(prefix string, fmtAndArgs ...interface{}) {
	msg := prefix
	if len(fmtAndArgs) > 0 {
		format, ok := fmtAndArgs[0].(string)
		if !ok {
			panic(errors.New("wrong fmt is given"))
		}
		msg += fmt.Sprintf(format, fmtAndArgs[1:]...)
	}
	panic(errors.New(msg))
}

func MustBeTrue(expr bool, fmtAndArgs ...interface{}) {
	if !expr {
		doPanic("expr must be true", fmtAndArgs...)
	}
}

func MustBeSuccess(err error, fmtAndArgs ...interface{}) {
	if err != nil {
		doPanic(fmt.Sprintf("error occurs: %s\n", err.Error()), fmtAndArgs...)
	}
}
