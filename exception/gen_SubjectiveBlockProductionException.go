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

var SubjectiveBlockProductionExceptionName = reflect.TypeOf(SubjectiveBlockProductionException{}).Name()

type SubjectiveBlockProductionException struct {
	_MiscException
	Elog log.Messages
}

func NewSubjectiveBlockProductionException(parent _MiscException, message log.Message) *SubjectiveBlockProductionException {
	return &SubjectiveBlockProductionException{parent, log.Messages{message}}
}

func (e SubjectiveBlockProductionException) Code() int64 {
	return 3100006
}

func (e SubjectiveBlockProductionException) Name() string {
	return SubjectiveBlockProductionExceptionName
}

func (e SubjectiveBlockProductionException) What() string {
	return "Subjective exception thrown during block production"
}

func (e *SubjectiveBlockProductionException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e SubjectiveBlockProductionException) GetLog() log.Messages {
	return e.Elog
}

func (e SubjectiveBlockProductionException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e SubjectiveBlockProductionException) DetailMessage() string {
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

func (e SubjectiveBlockProductionException) String() string {
	return e.DetailMessage()
}

func (e SubjectiveBlockProductionException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3100006,
		Name: SubjectiveBlockProductionExceptionName,
		What: "Subjective exception thrown during block production",
	}

	return json.Marshal(except)
}

func (e SubjectiveBlockProductionException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*SubjectiveBlockProductionException):
		callback(&e)
		return true
	case func(SubjectiveBlockProductionException):
		callback(e)
		return true
	default:
		return false
	}
}
