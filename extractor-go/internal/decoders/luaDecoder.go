package decoders

// statically linking liblua 5.1 x86 (because RO LUBs use this exact version)

//#cgo CFLAGS: -I${SRCDIR}/../../lua514/src
//#cgo LDFLAGS: -L${SRCDIR}/../../lua514/src -llua -lm
import "C"
import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/aarzilli/golua/lua"
	"github.com/guilherme-gm/ro-vis/extractor/internal/utils/stack"
)

type luaDecoder struct {
	L                *lua.State
	reencoder        StringReencoder
	path             *stack.Stack[string]
	notConsumedPaths map[string]bool
}

type LuaDecoderResult struct {
	NotConsumedPaths []string
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

func (decoder *luaDecoder) decodeSlice(slice reflect.Value, ctx luaDecContextInfo) {
	sliceType := slice.Type()
	sliceItemType := sliceType.Elem()

	newSlice := reflect.MakeSlice(sliceType, 0, 0)

	decoder.L.PushNil()
	for decoder.L.Next(-2) != 0 {
		sliceItem := reflect.New(sliceItemType).Elem()
		decoder.decode(sliceItem, newLuaDecContextInfo().setTableIndex(decoder.L.ToInteger(-2)))
		newSlice = reflect.Append(newSlice, sliceItem)

		decoder.L.Pop(1)
	}

	slice.Set(newSlice)
}

func (decoder *luaDecoder) decodeStruct(structObj reflect.Value, ctx luaDecContextInfo) {
	structType := structObj.Type()

	fieldList := make(map[string]bool)
	decoder.L.PushNil()
	for decoder.L.Next(-2) != 0 {
		switch decoder.L.Type(-2) {
		case lua.LUA_TSTRING:
			fieldName := decoder.L.ToString(-2)
			fieldList[fieldName] = true
		case lua.LUA_TNUMBER:
			fieldName := fmt.Sprintf("$$numeric:%d", decoder.L.ToInteger(-2))
			fieldList[fieldName] = true
		default:
			panic(fmt.Errorf("object key is not string. Found: %v", decoder.L.Type(-2)))
		}

		decoder.L.Pop(1)
	}

	for fldNum := range structType.NumField() {
		fieldType := structType.Field(fldNum)
		fieldValue := structObj.Field(fldNum)

		fieldName := fieldType.Name
		isKeyNumeric := false
		keyIndex := -1
		if alias := fieldType.Tag.Get("lua"); alias != "" {
			if alias == "@index" {
				if ctx.tableIndex == -1 {
					panic("Trying to get index of non-table")
				}

				fieldValue.SetInt(int64(ctx.tableIndex))
				continue
			}

			fieldName = alias
			if strings.HasPrefix(alias, "$$numeric:") {
				isKeyNumeric = true
				var err error
				keyIndex, err = strconv.Atoi(alias[10:])
				if err != nil {
					panic(err)
				}
			}
		}

		delete(fieldList, fieldName)

		if isKeyNumeric {
			decoder.L.PushInteger(int64(keyIndex))
			decoder.L.GetTable(-2)
		} else {
			decoder.L.GetField(-1, fieldName)
		}
		if decoder.L.IsNil(-1) {
			decoder.L.Pop(1)
			continue
		}

		decoder.path.Push(fieldName)
		decoder.decode(fieldValue, newLuaDecContextInfo())
		decoder.path.Pop()

		decoder.L.Pop(1)
	}

	for k := range fieldList {
		decoder.notConsumedPaths[strings.Join(decoder.path.ToSlice(), "/")+"/"+k] = true
	}
}

func (decoder *luaDecoder) decode(dataValue reflect.Value, ctx luaDecContextInfo) {
	dataType := dataValue.Type()
	dataKind := dataType.Kind()

	switch dataKind {
	case reflect.Slice:
		decoder.decodeSlice(dataValue, ctx)

	case reflect.Struct:
		decoder.decodeStruct(dataValue, ctx)

	case reflect.String:
		str := decoder.L.ToString(-1)
		dataValue.SetString(decoder.reencoder(str))

	case reflect.Int:
		val := decoder.L.ToInteger(-1)
		dataValue.SetInt(int64(val))

	case reflect.Bool:
		val := decoder.L.ToBoolean(-1)
		dataValue.SetBool(val)

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

func NewLuaDecoder(reencoder StringReencoder) *luaDecoder {
	decoder := &luaDecoder{
		L:                lua.NewState(),
		reencoder:        reencoder,
		path:             stack.NewStack[string](),
		notConsumedPaths: make(map[string]bool),
}

	decoder.L.OpenLibs()

	return decoder
}

func (decoder *luaDecoder) DecodeLuaTable(filePath string, tableName string, dst any) LuaDecoderResult {
	defer decoder.L.Close()

	err := decoder.L.DoFile(filePath)
	if err != nil {
		panic(err)
	}

	decoder.L.GetGlobal(tableName)
	decoder.path.Push(tableName)

	qv := reflect.ValueOf(dst)
	decoder.decode(qv.Elem(), newLuaDecContextInfo())

	decoder.path.Pop()

	var notConsumedPaths []string
	for k := range decoder.notConsumedPaths {
		notConsumedPaths = append(notConsumedPaths, k)
	}

	return LuaDecoderResult{
		NotConsumedPaths: notConsumedPaths,
	}
}

func DecodeLuaTable(filePath string, tableName string, dst any, reencoder StringReencoder) LuaDecoderResult {
	decoder := NewLuaDecoder(reencoder)
	return decoder.DecodeLuaTable(filePath, tableName, dst)
}
