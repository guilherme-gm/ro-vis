package decoders

import (
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

type myTableWithConsts struct {
	Id    int `lua:"@index"`
	Name  string
	Value int
}

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
