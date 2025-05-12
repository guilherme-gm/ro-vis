package decoders

// statically linking liblua 5.1 x86 (because RO LUBs use this exact version)

//#cgo CFLAGS: -I${SRCDIR}/../../../libs/liblua5.1/include
//#cgo LDFLAGS: -L${SRCDIR}/../../../libs/liblua5.1 -llua5.1 -lm
import "C"
import (
	"fmt"
	"reflect"

	"github.com/aarzilli/golua/lua"
)

type luaDecoder struct {
	L *lua.State
}

type luaDecContextInfo struct {
	tableIndex int
}

func newLuaDecContextInfo() luaDecContextInfo {
	return luaDecContextInfo{
		tableIndex: -1,
	}
}

func (c luaDecContextInfo) setTableIndex(index int) luaDecContextInfo {
	c.tableIndex = index
	return c
}

func (d *luaDecoder) decodeSlice(slice reflect.Value, ctx luaDecContextInfo) {
	sliceType := slice.Type()
	sliceItemType := sliceType.Elem()

	newSlice := reflect.MakeSlice(sliceType, 0, 0)

	d.L.PushNil()
	for d.L.Next(-2) != 0 {
		sliceItem := reflect.New(sliceItemType).Elem()
		d.decode(sliceItem, newLuaDecContextInfo().setTableIndex(d.L.ToInteger(-2)))
		newSlice = reflect.Append(newSlice, sliceItem)

		d.L.Pop(1)
	}

	slice.Set(newSlice)
}

func (d *luaDecoder) decodeStruct(structObj reflect.Value, ctx luaDecContextInfo) {
	structType := structObj.Type()

	fieldList := make(map[string]bool)
	d.L.PushNil()
	for d.L.Next(-2) != 0 {
		if d.L.Type(-2) != lua.LUA_TSTRING {
			panic("Object key is not string")
		}

		fieldName := d.L.ToString(-2)
		fieldList[fieldName] = true

		d.L.Pop(1)
	}

	for fldNum := range structType.NumField() {
		fieldType := structType.Field(fldNum)
		fieldValue := structObj.Field(fldNum)

		fieldName := fieldType.Name
		if alias := fieldType.Tag.Get("lua"); alias != "" {
			if alias == "@index" {
				if ctx.tableIndex == -1 {
					panic("Trying to get index of non-table")
				}

				fieldValue.SetInt(int64(ctx.tableIndex))
				continue
			}

			fieldName = alias
		}

		delete(fieldList, fieldName)

		d.L.GetField(-1, fieldName)
		if d.L.IsNil(-1) {
			d.L.Pop(1)
			continue
		}

		d.decode(fieldValue, newLuaDecContextInfo())

		d.L.Pop(1)
	}

	if len(fieldList) > 0 {
		fmt.Println("Not all keys were consumed.", fieldList)
		panic("Not all keys were consumed.")
	}
}

func (d *luaDecoder) decode(dataValue reflect.Value, ctx luaDecContextInfo) {
	dataType := dataValue.Type()
	dataKind := dataType.Kind()

	switch dataKind {
	case reflect.Slice:
		d.decodeSlice(dataValue, ctx)

	case reflect.Struct:
		d.decodeStruct(dataValue, ctx)

	case reflect.String:
		str := d.L.ToString(-1)
		dataValue.SetString(ConvertToUTF8(str))

	case reflect.Int:
		val := d.L.ToInteger(-1)
		dataValue.SetInt(int64(val))

	case reflect.Int8:
	case reflect.Uint8:
	case reflect.Int16:
	case reflect.Uint16:
	case reflect.Int32:
	case reflect.Uint32:
	case reflect.Int64:
	case reflect.Uint64:
		panic("LuaDecoder doesn't handle sized int fields properly. use int. Found: " + dataValue.String())

	default:
		panic("decode default - " + dataValue.String())
	}
}

func DecodeLuaTable(filePath string, tableName string, dst any) {
	decoder := luaDecoder{
		L: lua.NewState(),
	}
	decoder.L.OpenLibs()
	defer decoder.L.Close()

	err := decoder.L.DoFile(filePath)
	if err != nil {
		panic(err)
	}

	decoder.L.GetGlobal(tableName)
	qv := reflect.ValueOf(dst)
	decoder.decode(qv.Elem(), newLuaDecContextInfo())
}
