// Code generated by gotemplate. DO NOT EDIT.

package exception

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/eosspark/eos-go/log"
)

// template type Exception(PARENT,CODE,WHAT)

var UnaccessibleApiName = reflect.TypeOf(UnaccessibleApi{}).Name()

type UnaccessibleApi struct {
	_ActionValidateException
	Elog log.Messages
}

func NewUnaccessibleApi(parent _ActionValidateException, message log.Message) *UnaccessibleApi {
	return &UnaccessibleApi{parent, log.Messages{message}}
}

func (e UnaccessibleApi) Code() int64 {
	return 3050007
}

func (e UnaccessibleApi) Name() string {
	return UnaccessibleApiName
}

func (e UnaccessibleApi) What() string {
	return "Attempt to use unaccessible API"
}

func (e *UnaccessibleApi) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e UnaccessibleApi) GetLog() log.Messages {
	return e.Elog
}

func (e UnaccessibleApi) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e UnaccessibleApi) DetailMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(e.Code())))
	buffer.WriteString(" ")
	buffer.WriteString(e.Name())
	buffer.WriteString(": ")
	buffer.WriteString(e.What())
	buffer.WriteString("\n")
	for _, l := range e.Elog {
		buffer.WriteString("[")
		buffer.WriteString(l.GetMessage())
		buffer.WriteString("] ")
		buffer.WriteString(l.GetContext().String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (e UnaccessibleApi) String() string {
	return e.DetailMessage()
}

func (e UnaccessibleApi) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3050007,
		Name: UnaccessibleApiName,
		What: "Attempt to use unaccessible API",
	}

	return json.Marshal(except)
}

func (e UnaccessibleApi) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*UnaccessibleApi):
		callback(&e)
		return true
	case func(UnaccessibleApi):
		callback(e)
		return true
	default:
		return false
	}
}
