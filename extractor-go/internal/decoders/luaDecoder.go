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

/**
 * General purpose Lua file decoder.
 *
 * It executes Lua files and is able to extract Lua tables into Go structs, by scanning Lua VM's data.
 *
 * The provided struct must represent the table structure, both in field names and use the right types.
 *
 * The following tags are supported:
 * - `lua:"<name>"`: used to provide the corresponding name in Lua table (when it does not match the struct field name)
 *
 * - `lua:"@index"`: used to mark that the field should be populated with the table item "index"
 *    For example, given this table: tbl = { [1001] = { Value = "A" } }
 *    The struct field marked as `lua:"@index"` will be populated with 1001
 *
 * - `lua:"$$numeric:<pos>"`: used to mark that the field should be populated with the value at index <pos>
 *    For example, given this table: tbl = { [1001] = { "A", "B" } }
 *    The struct field marked as `lua:"$$numeric:2"` will be populated with "B"
 *
 * - `lua:"@sliceValue"`: used to mark that this field should expand the current table item as an array, where each item is a struct
 *    For example, given this table:
 *    tbl = {
 *       [1001] = {
 *          Values = {
 *             [100] = { -- !!! This is an array of structs
 *                { FieldA = "DataA", FieldB = "DataB" },
 *                { FieldA = "DataA", FieldB = "DataB" },
 *             },
 *             [101] = { }
 *          }
 *       }
 *    }
 *
 *    Could be represented as:
 *    type MyTable struct {
 *       Id    int `lua:"@index"` -- 1001
 *       Values []struct {
 *          Id    int `lua:"@index"` -- 100, 101
 *          Data []struct {
 *             FieldA string
 *             FieldB string
 *          } `lua:"@sliceValue"` -- !!! Expand the values of 100/101 as elements of this struct.
 *       }
 *    }
 *
 * - `lua:"@plain"`: used to mark that a slice actually represents a dynamic list of key/value pairs.
 *    Usually used together with `lua:"@index"` and `lua:"@plainValue"` to read data.
 *
 *   For example:
 *   tbl = {
 *      Key1 = 1,
 *      Key2 = 2,
 *   }
 *
 *   Could be represented as:
 *   type MyTable struct {
 *      Values []struct {
 *         Key string `lua:"@index"` -- "Key1", "Key2"
 *         Value int `lua:"@plainValue"` -- 1, 2
 *      } `lua:"@plain"`
 *   }
 */

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
	tableIndexType lua.LuaValType
	tableIndexStr  string
	tableIndex     int
	isPlain        bool
}

func newLuaDecContextInfo(previousCtx *luaDecContextInfo) luaDecContextInfo {
	newCtx := luaDecContextInfo{
		tableIndexType: lua.LUA_TNIL,
	}

	if previousCtx != nil {
		newCtx.isPlain = previousCtx.isPlain
	}

	return newCtx
}

func (c luaDecContextInfo) setTableIndex(L *lua.State, index int) luaDecContextInfo {
	switch L.Type(index) {
	case lua.LUA_TSTRING:
		c.tableIndexStr = L.ToString(index)
		c.tableIndexType = lua.LUA_TSTRING
	case lua.LUA_TNUMBER:
		c.tableIndex = L.ToInteger(index)
		c.tableIndexType = lua.LUA_TNUMBER
	default:
		panic(fmt.Errorf("unexpected table index type: %v", L.Type(index)))
	}

	return c
}

func (decoder *luaDecoder) decodeSlice(slice reflect.Value, ctx luaDecContextInfo) {
	sliceType := slice.Type()
	sliceItemType := sliceType.Elem()

	newSlice := reflect.MakeSlice(sliceType, 0, 0)

	// fmt.Println("Decoding slice: " + sliceType.String())
	decoder.L.PushNil()
	for decoder.L.Next(-2) != 0 {
		newCtx := newLuaDecContextInfo(&ctx).setTableIndex(decoder.L, -2)
		// __newindex is a metatable field, we don't mind them
		if newCtx.tableIndexType == lua.LUA_TSTRING && newCtx.tableIndexStr == "__newindex" {
			decoder.L.Pop(1)
			continue
		}

		sliceItem := reflect.New(sliceItemType).Elem()
		decoder.decode(sliceItem, newCtx)
		newSlice = reflect.Append(newSlice, sliceItem)

		decoder.L.Pop(1)
	}

	slice.Set(newSlice)
}

func (decoder *luaDecoder) decodeStruct(structObj reflect.Value, ctx luaDecContextInfo) {
	structType := structObj.Type()
	// fmt.Println("Decoding struct: " + structType.String())
	// fmt.Println("-- CTX: int: " + fmt.Sprintf("%d", ctx.tableIndex))
	// fmt.Println("-- CTX: string: " + ctx.tableIndexStr)
	// fmt.Println("-- CTX: plain: " + fmt.Sprintf("%t", ctx.isPlain))

	fieldList := make(map[string]bool)
	// Plain structures must not be enumerated, as they are not tables
	if !ctx.isPlain {
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
	}

	for fldNum := range structType.NumField() {
		// fmt.Println("Decoding field: " + structType.Field(fldNum).Name)

		fieldType := structType.Field(fldNum)
		fieldValue := structObj.Field(fldNum)

		fieldName := fieldType.Name
		isKeyNumeric := false
		isPlain := false
		keyIndex := -1
		if alias := fieldType.Tag.Get("lua"); alias != "" {
			if alias == "@index" {
				if ctx.tableIndex == -1 {
					panic("Trying to get index of non-table")
				}

				switch fieldType.Type.Kind() {
				case reflect.Int:
					if ctx.tableIndexType != lua.LUA_TNUMBER {
						panic(fmt.Sprintf("Incompatible type for @index. Go Type: %s | Lua TYPE: %d | String value: %s", fieldType.Type.String(), ctx.tableIndexType, ctx.tableIndexStr))
					}
					fieldValue.SetInt(int64(ctx.tableIndex))

				case reflect.String:
					if ctx.tableIndexType != lua.LUA_TSTRING {
						panic(fmt.Sprintf("Incompatible type for @index. Go Type: %s | Lua TYPE: %d", fieldType.Type.String(), ctx.tableIndexType))
					}
					fieldValue.SetString(ctx.tableIndexStr)

				default:
					panic("Unsupported type for @index: " + fieldType.Type.String())
				}
				continue
			}

			if alias == "@plain" {
				isPlain = true
			}

			if alias == "@plainValue" {
				decoder.decode(fieldValue, newLuaDecContextInfo(&ctx))
				continue
			}

			if alias == "@sliceValue" {
				decoder.decodeSlice(fieldValue, newLuaDecContextInfo(&ctx))
				// remove all $$numeric from fieldList, as they were likely handled by this slice
				for k := range fieldList {
					if strings.HasPrefix(k, "$$numeric:") {
						delete(fieldList, k)
					}
				}
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
		} else if isPlain {
			/* do nothing */
		} else {
			decoder.L.GetField(-1, fieldName)
		}

		if decoder.L.IsNil(-1) {
			decoder.L.Pop(1)
			continue
		}

		decoder.path.Push(fieldName)
		newCtx := newLuaDecContextInfo(&ctx)
		newCtx.isPlain = true
		decoder.decode(fieldValue, newCtx)
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

	// fmt.Println("Decoding type: " + dataType.String())
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

// creates a table with string keys and int values
// this allows you to add custom tables that the parser will use when decoding a file
func (decoder *luaDecoder) CreateIntTable(tableName string, values map[string]int) {
	decoder.L.CreateTable(0, len(values)) // pushes the new table into stack

	for k, v := range values {
		decoder.L.PushString(k)         // Table key
		decoder.L.PushInteger(int64(v)) // Table value
		decoder.L.SetTable(-3)          // Makes t[key] = value and pops key and value
	}

	decoder.L.SetGlobal(tableName) // Pops the table from stack and set it to a variable
}

// creates a variable with int value
// this allows you to add custom variables that the parser will use when decoding a file
func (decoder *luaDecoder) CreateIntVar(varName string, value int) {
	decoder.L.PushInteger(int64(value)) // pushes the new table into stack
	decoder.L.SetGlobal(varName)        // Pops the table from stack and set it to a variable
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
	decoder.decode(qv.Elem(), newLuaDecContextInfo(nil))

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
