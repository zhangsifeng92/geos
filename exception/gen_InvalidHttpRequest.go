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

var InvalidHttpRequestName = reflect.TypeOf(InvalidHttpRequest{}).Name()

type InvalidHttpRequest struct {
	_HttpException
	Elog log.Messages
}

func NewInvalidHttpRequest(parent _HttpException, message log.Message) *InvalidHttpRequest {
	return &InvalidHttpRequest{parent, log.Messages{message}}
}

func (e InvalidHttpRequest) Code() int64 {
	return 3200006
}

func (e InvalidHttpRequest) Name() string {
	return InvalidHttpRequestName
}

func (e InvalidHttpRequest) What() string {
	return "invalid http request"
}

func (e *InvalidHttpRequest) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e InvalidHttpRequest) GetLog() log.Messages {
	return e.Elog
}

func (e InvalidHttpRequest) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e InvalidHttpRequest) DetailMessage() string {
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
		buffer.WriteString("]")
		buffer.WriteString("\n")
		buffer.WriteString(l.GetContext().String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (e InvalidHttpRequest) String() string {
	return e.DetailMessage()
}

func (e InvalidHttpRequest) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3200006,
		Name: InvalidHttpRequestName,
		What: "invalid http request",
	}

	return json.Marshal(except)
}

func (e InvalidHttpRequest) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*InvalidHttpRequest):
		callback(&e)
		return true
	case func(InvalidHttpRequest):
		callback(e)
		return true
	default:
		return false
	}
}