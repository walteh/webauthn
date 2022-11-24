package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type Error struct {
	Code_ uint8
	// Short name for the type of error that has occurred
	Type_ string `json:"type"`
	// Additional details about the error
	Details_ string `json:"error"`
	// Information to help debug the error
	DevInfo_ string `json:"debug"`

	Caller_ string `json:"caller"`

	KV_ map[string]interface{} `json:"data"`

	Root error `json:"root"`
}

func NewError(code uint8) *Error {
	return &Error{Code_: code}
}

func (err *Error) WithCaller() *Error {
	// add line number and file name to caller
	_, file, no, ok := runtime.Caller(1)
	if ok {
		err.Caller_ = fmt.Sprintf("%s:%d", file, no)
	} else {
		err.Caller_ = "unknown"
	}

	return err

}

func (err *Error) Copy() *Error {
	newErr := *err
	return &newErr
}

func (err *Error) Type() string {
	return err.Type_
}

func (err *Error) Message() string {
	return err.Details_
}

func (err *Error) Caller() string {
	// format the caller so it is just the last two parts of the path
	// and the line number

	// split the caller string by the / character
	split := strings.Split(err.Caller_, "/")

	// if the split is less than 2, just return the caller
	if len(split) < 2 {
		return err.Caller_
	}

	// return the last two parts of the split
	return fmt.Sprintf("%s/%s", split[len(split)-2], split[len(split)-1])
}

func (err *Error) DevInfo() string {
	return err.DevInfo_
}

func (err *Error) Roots() []error {
	var wrk error
	wrk = err.Copy()
	ok := true
	var s []error
	for ok {
		var er *Error
		s = append(s, wrk)
		if er, ok = wrk.(*Error); ok {
			if er.Root != nil {
				wrk = er.Root
			} else {
				ok = false
			}
		} else {
			ok = false
		}
	}

	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s

}

func (err *Error) KV() map[string]interface{} {
	return err.KV_
}

func (err *Error) Code() string {
	// format code as hex
	return "0x" + fmt.Sprintf("%x", err.Code_)
}

func (err *Error) Error() string {
	// format the error as a pretty string
	// for logging and debugging
	// var str string
	str := fmt.Sprintf("%s - %s - %s - %s - %s - %s", err.Code(), err.Type(), err.Caller(), err.Message(), err.DevInfo(), err.KV())

	// // add the root error if it exists
	// if err.Root != nil {
	// 	str = fmt.Sprintf("%s\n\t\t %s", str, err.Root.Error())
	// }
	return string(str)
}

func (passedError *Error) WithMessage(details string) *Error {
	err := passedError.Copy()
	err.Details_ = details
	return err
}

func (passedError *Error) WithRoot(root error) *Error {
	err := passedError.Copy()
	err.Root = root
	return err
}

func (passedError *Error) WithInfo(info string) *Error {
	err := passedError.Copy()
	err.DevInfo_ = info
	return err
}

func (passedError *Error) WithType(t string) *Error {
	err := passedError.Copy()
	err.Type_ = t
	return err
}

func (err *Error) WithKV(str string, inter interface{}) *Error {
	if err.KV_ == nil {
		err.KV_ = make(map[string]interface{})
	}
	err.KV_[str] = inter
	return err
}
