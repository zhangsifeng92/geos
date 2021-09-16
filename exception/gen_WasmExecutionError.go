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

var WasmExecutionErrorName = reflect.TypeOf(WasmExecutionError{}).Name()

type WasmExecutionError struct {
	_WasmException
	Elog log.Messages
}

func NewWasmExecutionError(parent _WasmException, message log.Message) *WasmExecutionError {
	return &WasmExecutionError{parent, log.Messages{message}}
}

func (e WasmExecutionError) Code() int64 {
	return 3070002
}

func (e WasmExecutionError) Name() string {
	return WasmExecutionErrorName
}

func (e WasmExecutionError) What() string {
	return "Runtime Error Processing WASM"
}

func (e *WasmExecutionError) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e WasmExecutionError) GetLog() log.Messages {
	return e.Elog
}

func (e WasmExecutionError) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e WasmExecutionError) DetailMessage() string {
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

func (e WasmExecutionError) String() string {
	return e.DetailMessage()
}

func (e WasmExecutionError) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3070002,
		Name: WasmExecutionErrorName,
		What: "Runtime Error Processing WASM",
	}

	return json.Marshal(except)
}

func (e WasmExecutionError) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*WasmExecutionError):
		callback(&e)
		return true
	case func(WasmExecutionError):
		callback(e)
		return true
	default:
		return false
	}
}
