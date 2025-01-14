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

var InvalidHttpResponseName = reflect.TypeOf(InvalidHttpResponse{}).Name()

type InvalidHttpResponse struct {
	_HttpException
	Elog log.Messages
}

func NewInvalidHttpResponse(parent _HttpException, message log.Message) *InvalidHttpResponse {
	return &InvalidHttpResponse{parent, log.Messages{message}}
}

func (e InvalidHttpResponse) Code() int64 {
	return 3200002
}

func (e InvalidHttpResponse) Name() string {
	return InvalidHttpResponseName
}

func (e InvalidHttpResponse) What() string {
	return "invalid http response"
}

func (e *InvalidHttpResponse) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e InvalidHttpResponse) GetLog() log.Messages {
	return e.Elog
}

func (e InvalidHttpResponse) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e InvalidHttpResponse) DetailMessage() string {
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

func (e InvalidHttpResponse) String() string {
	return e.DetailMessage()
}

func (e InvalidHttpResponse) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3200002,
		Name: InvalidHttpResponseName,
		What: "invalid http response",
	}

	return json.Marshal(except)
}

func (e InvalidHttpResponse) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*InvalidHttpResponse):
		callback(&e)
		return true
	case func(InvalidHttpResponse):
		callback(e)
		return true
	default:
		return false
	}
}
