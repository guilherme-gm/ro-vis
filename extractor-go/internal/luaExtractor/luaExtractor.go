package luaExtractor

// statically linking liblua 5.1 x86 (because RO LUBs use this exact version)

//#cgo CFLAGS: -I${SRCDIR}/../../../libs/liblua5.1/include
//#cgo LDFLAGS: -L${SRCDIR}/../../../libs/liblua5.1 -llua5.1 -lm
import "C"
import (
	"bytes"
	"io"
	"reflect"

	"golang.org/x/net/html/charset"

	"github.com/aarzilli/golua/lua"
)

type QuestV1 struct {
	QuestId     int `lua:"@index"`
	Title       string
	Description []string
	Summary     string
}

func convertToUTF8(str string) string {
	strBytes := []byte(str)
	byteReader := bytes.NewReader(strBytes)
	reader, _ := charset.NewReaderLabel("euc-kr", byteReader)
	strBytes, _ = io.ReadAll(reader)
	return string(strBytes)
}

func decodeSlice(L *lua.State, slice reflect.Value) {
	sliceType := slice.Type()
	sliceItemType := sliceType.Elem()

	newSlice := reflect.MakeSlice(sliceType, 0, 0)

	L.PushNil()
	for L.Next(-2) != 0 {
		sliceItem := reflect.New(sliceItemType).Elem()

		decode(L, sliceItem)
		newSlice = reflect.Append(slice, sliceItem)

		L.Pop(1)
	}

	slice.Set(newSlice)
}

func decodeStruct(L *lua.State, structObj reflect.Value) {
	structType := structObj.Type()
	for fldNum := range structType.NumField() {
		fieldType := structType.Field(fldNum)
		fieldValue := structObj.Field(fldNum)

		if alias := fieldType.Tag.Get("lua"); alias != "" {
			continue // @TODO:
		}

		L.GetField(-1, fieldType.Name)
		if L.IsNil(-1) {
			L.Pop(1)
			continue
		}

		decode(L, fieldValue)

		L.Pop(1)
	}
}

func decode(L *lua.State, dataValue reflect.Value) {
	dataType := dataValue.Type()
	dataKind := dataType.Kind()

	switch dataKind {
	case reflect.Slice:
		decodeSlice(L, dataValue)

	case reflect.Struct:
		decodeStruct(L, dataValue)

	case reflect.String:
		str := L.ToString(-1)
		dataValue.SetString(convertToUTF8(str))

	default:
		panic("decode default")
	}
}

func Decode(filePath string, tableName string, dst any) {
	L := lua.NewState()
	L.OpenLibs()
	defer L.Close()

	err := L.DoFile(filePath)
	if err != nil {
		panic(err)
	}

	L.GetGlobal(tableName)
	qv := reflect.ValueOf(dst)
	decode(L, qv.Elem())
}
