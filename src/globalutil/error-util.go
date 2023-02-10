package globalutil

import (
	"errors"
	"fmt"
	"runtime"
)

type ErrorShower interface {
	ErrorShow() string
}

type ErrorLocer interface {
	ErrorLoc() string
}

type FormatedError struct {
	ErrText string `json:"error"`
}

type errDetail struct {
	loc  string // Location of error
	show string // Error message to show to user
	err  error  // Error to log
}

func (e *errDetail) Error() string {
	return fmt.Sprintf("ErrorDetail: show=%q, loc=%q, err: %s", e.show, e.loc, e.err.Error())
}

func (e *errDetail) ErrorLoc() string {
	return e.loc
}

func (e *errDetail) ErrorShow() string {
	return e.show
}

func Errorf(cause error, fmtstr string, args ...interface{}) error {
	showStr := fmt.Sprintf(fmtstr, args...)

	_, file, line, _ := runtime.Caller(1)

	if cause == nil {
		return &errDetail{
			loc:  fmt.Sprintf("%s:%d", file, line),
			show: showStr,
			err:  errors.New(showStr),
		}
	}

	errDet, ok := cause.(*errDetail)
	if !ok {
		return &errDetail{
			loc:  fmt.Sprintf("%s:%d", file, line),
			show: showStr,
			err:  cause,
		}
	}

	if errDet.show == "" {
		errDet.show = showStr
	}

	if errDet.loc != "" {
		errDet.loc = fmt.Sprintf("%s:%d :: %s", file, line, errDet.loc)
	}

	return errDet
}
