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

var OverlappingMemoryErrorName = reflect.TypeOf(OverlappingMemoryError{}).Name()

type OverlappingMemoryError struct {
	_WasmException
	Elog log.Messages
}

func NewOverlappingMemoryError(parent _WasmException, message log.Message) *OverlappingMemoryError {
	return &OverlappingMemoryError{parent, log.Messages{message}}
}

func (e OverlappingMemoryError) Code() int64 {
	return 3070004
}

func (e OverlappingMemoryError) Name() string {
	return OverlappingMemoryErrorName
}

func (e OverlappingMemoryError) What() string {
	return "memcpy with overlapping memory"
}

func (e *OverlappingMemoryError) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e OverlappingMemoryError) GetLog() log.Messages {
	return e.Elog
}

func (e OverlappingMemoryError) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e OverlappingMemoryError) DetailMessage() string {
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

func (e OverlappingMemoryError) String() string {
	return e.DetailMessage()
}

func (e OverlappingMemoryError) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3070004,
		Name: OverlappingMemoryErrorName,
		What: "memcpy with overlapping memory",
	}

	return json.Marshal(except)
}

func (e OverlappingMemoryError) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*OverlappingMemoryError):
		callback(&e)
		return true
	case func(OverlappingMemoryError):
		callback(e)
		return true
	default:
		return false
	}
}
