package decoders

import (
	"slices"
	"testing"

	testfiles "github.com/guilherme-gm/ro-vis/extractor/testfiles"
	"github.com/stretchr/testify/assert"
)

type myTable struct {
	Name  string
	Value int
}

type myTableWithIds struct {
	Id    int `lua:"@index"`
	Name  string
	Value int
}

type myTableIntArray struct {
	Id     int `lua:"@index"`
	Values []int
}

type myTableWithSliceValue struct {
	Id   int `lua:"@index"`
	Data []struct {
		FieldA string
		FieldB string
	} `lua:"@sliceValue"`
}

type myTableIndexedValue struct {
	Id        int    `lua:"@index"`
	FirstVal  string `lua:"$$numeric:1"`
	SecondVal int    `lua:"$$numeric:2"`
}

type myTableWithSlices struct {
	Id    int `lua:"@index"`
	Name  string
	Value []myTableWithSliceValue
}

type myTableWithConsts struct {
	Id    int `lua:"@index"`
	Name  string
	Value int
}

// Tests decoding a simple table into a struct
func TestLuaInstanceDecodeTableArray(t *testing.T) {
	var dst []myTable

	decoder := NewLuaDecoder(ConvertNoop)
	decoder.DecodeLuaTable(testfiles.GetFilePath("lua_tables.lua"), "MY_TABLE", &dst)

	assert.Equal(t, 2, len(dst))
	assert.Equal(t, "Test", dst[0].Name)
	assert.Equal(t, 1, dst[0].Value)
	assert.Equal(t, "Test2", dst[1].Name)
	assert.Equal(t, 2, dst[1].Value)
}

// Tests decoding tables indexes as fields with `lua:"@index"`
func TestLuaInstanceDecodeTableWithIds(t *testing.T) {
	var dst []myTableWithIds

	decoder := NewLuaDecoder(ConvertNoop)
	decoder.DecodeLuaTable(testfiles.GetFilePath("lua_tables.lua"), "MY_TABLE_WITH_IDS", &dst)

	assert.Equal(t, 2, len(dst))

	assert.Equal(t, 1005, dst[0].Id)
	assert.Equal(t, "Test", dst[0].Name)
	assert.Equal(t, 1, dst[0].Value)

	assert.Equal(t, 1006, dst[1].Id)
	assert.Equal(t, "Test2", dst[1].Name)
	assert.Equal(t, 2, dst[1].Value)
}

// Tests decoding a table entry where values are simply indexed as an "array"
func TestLuaInstanceDecodeTableIntArray(t *testing.T) {
	var dst []myTableIntArray

	decoder := NewLuaDecoder(ConvertNoop)
	decoder.DecodeLuaTable(testfiles.GetFilePath("lua_tables.lua"), "MY_TABLE_INT_ARRAY", &dst)

	assert.Equal(t, 3, len(dst))

	index := slices.IndexFunc(dst, func(item myTableIntArray) bool {
		return item.Id == 1005
	})
	assert.Equal(t, 1005, dst[index].Id)
	assert.Equal(t, []int{1, 2, 3}, dst[index].Values)

	index = slices.IndexFunc(dst, func(item myTableIntArray) bool {
		return item.Id == 1006
	})
	assert.Equal(t, 1006, dst[index].Id)
	assert.Equal(t, []int{}, dst[index].Values)

	index = slices.IndexFunc(dst, func(item myTableIntArray) bool {
		return item.Id == 1007
	})
	assert.Equal(t, 1007, dst[index].Id)
	assert.Nil(t, dst[index].Values)
}

// Tests decoding a table entry where values are simply indexed as an "array", using $$numeric
func TestLuaInstanceDecodeTableWithIndexedValue(t *testing.T) {
	var dst []myTableIndexedValue

	decoder := NewLuaDecoder(ConvertNoop)
	decoder.DecodeLuaTable(testfiles.GetFilePath("lua_tables.lua"), "MY_TABLE_INDEXED", &dst)

	assert.Equal(t, 2, len(dst))

	assert.Equal(t, 1005, dst[0].Id)
	assert.Equal(t, "Value1", dst[0].FirstVal)
	assert.Equal(t, 10, dst[0].SecondVal)

	assert.Equal(t, 1006, dst[1].Id)
	assert.Equal(t, "Value2", dst[1].FirstVal)
	assert.Equal(t, 20, dst[1].SecondVal)
}

// Tests decoding inner slices with `lua:"@sliceValue"`
func TestLuaInstanceDecodeTableWithSlices(t *testing.T) {
	var dst []myTableWithSlices

	decoder := NewLuaDecoder(ConvertNoop)
	decoder.DecodeLuaTable(testfiles.GetFilePath("lua_tables.lua"), "MY_TABLE_SLICES", &dst)

	assert.Equal(t, 2, len(dst))

	assert.Equal(t, 1005, dst[0].Id)
	assert.Equal(t, "Test", dst[0].Name)
	assert.Equal(t, 2, len(dst[0].Value))
	assert.Equal(t, 100, dst[0].Value[0].Id)
	assert.Equal(t, "DataA", dst[0].Value[0].Data[0].FieldA)
	assert.Equal(t, "DataB", dst[0].Value[0].Data[0].FieldB)

	assert.Equal(t, 101, dst[0].Value[1].Id)
	assert.Equal(t, "DataC", dst[0].Value[1].Data[0].FieldA)
	assert.Equal(t, "DataD", dst[0].Value[1].Data[0].FieldB)

	assert.Equal(t, 1006, dst[1].Id)
	assert.Equal(t, "Test2", dst[1].Name)
	assert.Equal(t, 0, len(dst[1].Value))
}

// Tests decoding a table with extra constants provided via code
func TestLuaInstanceDecodeTableWithExtraConsts(t *testing.T) {
	var dst []myTableWithConsts

	myDecoder := NewLuaDecoder(ConvertNoop)
	myDecoder.CreateIntTable("CON1", map[string]int{
		"Key1": 1005,
		"Key2": 1006,
	})
	myDecoder.CreateIntTable("CON2", map[string]int{
		"Val1": 10000,
		"Val2": 20000,
	})

	myDecoder.DecodeLuaTable(testfiles.GetFilePath("lua_table_with_consts.lua"), "MY_TABLE_WITH_CONSTS", &dst)

	assert.Equal(t, 2, len(dst))

	assert.Equal(t, 1005, dst[0].Id)
	assert.Equal(t, "Test", dst[0].Name)
	assert.Equal(t, 10000, dst[0].Value)

	assert.Equal(t, 1006, dst[1].Id)
	assert.Equal(t, "Test2", dst[1].Name)
	assert.Equal(t, 20000, dst[1].Value)
}

func TestLuaInstanceDecodeConstantsTable(t *testing.T) {
	var dstTable struct {
		Values []struct {
			Id    string `lua:"@index"`
			Value int    `lua:"@plainValue"`
		} `lua:"@plain"`
	}

	decoder := NewLuaDecoder(ConvertNoop)
	decoder.DecodeLuaTable(testfiles.GetFilePath("lua_tables.lua"), "MY_TABLE_INT", &dstTable)

	assert.Equal(t, 2, len(dstTable.Values))
	assert.Equal(t, "Key1", dstTable.Values[0].Id)
	assert.Equal(t, 1, dstTable.Values[0].Value)
	assert.Equal(t, "Key2", dstTable.Values[1].Id)
	assert.Equal(t, 2, dstTable.Values[1].Value)
}

func TestLuaInstanceDecodeTableWithMetaindex(t *testing.T) {
	var dstTable struct {
		Values []struct {
			Id    string `lua:"@index"`
			Value int    `lua:"@plainValue"`
		} `lua:"@plain"`
	}

	decoder := NewLuaDecoder(ConvertNoop)
	decoder.DecodeLuaTable(testfiles.GetFilePath("lua_tables.lua"), "MY_TABLE_INT_METATABLE", &dstTable)

	assert.Equal(t, 2, len(dstTable.Values))
	assert.Equal(t, "Key1", dstTable.Values[0].Id)
	assert.Equal(t, 1, dstTable.Values[0].Value)
	assert.Equal(t, "Key2", dstTable.Values[1].Id)
	assert.Equal(t, 2, dstTable.Values[1].Value)
}
