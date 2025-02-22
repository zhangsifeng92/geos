// Code generated by gotemplate. DO NOT EDIT.

package exception

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/zhangsifeng92/geos/log"
)

// template type Exception(PARENT,CODE,WHAT)

var ReversibleGuardExceptionName = reflect.TypeOf(ReversibleGuardException{}).Name()

type ReversibleGuardException struct {
	_GuardException
	Elog log.Messages
}

func NewReversibleGuardException(parent _GuardException, message log.Message) *ReversibleGuardException {
	return &ReversibleGuardException{parent, log.Messages{message}}
}

func (e ReversibleGuardException) Code() int64 {
	return 3060102
}

func (e ReversibleGuardException) Name() string {
	return ReversibleGuardExceptionName
}

func (e ReversibleGuardException) What() string {
	return "Reversible block log usage is at unsafe levels"
}

func (e *ReversibleGuardException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e ReversibleGuardException) GetLog() log.Messages {
	return e.Elog
}

func (e ReversibleGuardException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e ReversibleGuardException) DetailMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(e.Code())))
	buffer.WriteByte(' ')
	buffer.WriteString(e.Name())
	buffer.Write([]byte{':', ' '})
	buffer.WriteString(e.What())
	buffer.WriteByte('\n')
	for _, l := range e.Elog {
		buffer.WriteByte('[')
		buffer.WriteString(l.GetMessage())
		buffer.Write([]byte{']', ' '})
		buffer.WriteString(l.GetContext().String())
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

func (e ReversibleGuardException) String() string {
	return e.DetailMessage()
}

func (e ReversibleGuardException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3060102,
		Name: ReversibleGuardExceptionName,
		What: "Reversible block log usage is at unsafe levels",
	}

	return json.Marshal(except)
}

func (e ReversibleGuardException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*ReversibleGuardException):
		callback(&e)
		return true
	case func(ReversibleGuardException):
		callback(e)
		return true
	default:
		return false
	}
}
