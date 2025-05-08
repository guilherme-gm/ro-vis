package luaExtractor

// statically linking liblua 5.1 x86 (because RO LUBs use this exact version)

//#cgo CFLAGS: -I${SRCDIR}/../../../libs/liblua5.1/include
//#cgo LDFLAGS: -L${SRCDIR}/../../../libs/liblua5.1 -llua5.1 -lm
import "C"
import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"golang.org/x/net/html/charset"

	"github.com/aarzilli/golua/lua"
)

type contextInfo struct {
	tableIndex int
}

func newContextInfo() contextInfo {
	return contextInfo{
		tableIndex: -1,
	}
}

func (c contextInfo) setTableIndex(index int) contextInfo {
	c.tableIndex = index
	return c
}

func convertToUTF8(str string) string {
	strBytes := []byte(str)
	byteReader := bytes.NewReader(strBytes)
	reader, _ := charset.NewReaderLabel("euc-kr", byteReader)
	strBytes, _ = io.ReadAll(reader)
	return string(strBytes)
}

func decodeSlice(L *lua.State, slice reflect.Value, ctx contextInfo) {
	sliceType := slice.Type()
	sliceItemType := sliceType.Elem()

	newSlice := reflect.MakeSlice(sliceType, 0, 0)

	L.PushNil()
	for L.Next(-2) != 0 {
		sliceItem := reflect.New(sliceItemType).Elem()
		decode(L, sliceItem, newContextInfo().setTableIndex(L.ToInteger(-2)))
		newSlice = reflect.Append(newSlice, sliceItem)

		L.Pop(1)
	}

	slice.Set(newSlice)
}

func decodeStruct(L *lua.State, structObj reflect.Value, ctx contextInfo) {
	structType := structObj.Type()

	fieldList := make(map[string]bool)
	L.PushNil()
	for L.Next(-2) != 0 {
		if L.Type(-2) != lua.LUA_TSTRING {
			panic("Object key is not string")
		}

		fieldName := L.ToString(-2)
		fieldList[fieldName] = true

		L.Pop(1)
	}

	for fldNum := range structType.NumField() {
		fieldType := structType.Field(fldNum)
		fieldValue := structObj.Field(fldNum)

		delete(fieldList, fieldType.Name)

		if alias := fieldType.Tag.Get("lua"); alias != "" {
			if alias == "@index" {
				if ctx.tableIndex == -1 {
					panic("Trying to get index of non-table")
				}

				fieldValue.SetInt(int64(ctx.tableIndex))
				continue
			}

			panic("Invalid lua alias: " + alias)
		}

		L.GetField(-1, fieldType.Name)
		if L.IsNil(-1) {
			L.Pop(1)
			continue
		}

		decode(L, fieldValue, newContextInfo())

		L.Pop(1)
	}

	if len(fieldList) > 0 {
		fmt.Println("Not all keys were consumed.", fieldList)
		panic("Not all keys were consumed.")
	}
}

func decode(L *lua.State, dataValue reflect.Value, ctx contextInfo) {
	dataType := dataValue.Type()
	dataKind := dataType.Kind()

	switch dataKind {
	case reflect.Slice:
		decodeSlice(L, dataValue, ctx)

	case reflect.Struct:
		decodeStruct(L, dataValue, ctx)

	case reflect.String:
		str := L.ToString(-1)
		dataValue.SetString(convertToUTF8(str))

	case reflect.Int8:
	case reflect.Uint8:
	case reflect.Int16:
	case reflect.Uint16:
	case reflect.Int32:
	case reflect.Uint32:
	case reflect.Int64:
	case reflect.Uint64:
		val := L.ToInteger(-1)
		dataValue.SetInt(int64(val))

	default:
		panic("decode default - " + dataValue.String())
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
	decode(L, qv.Elem(), newContextInfo())
}
