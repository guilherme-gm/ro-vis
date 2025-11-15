MY_TABLE = {
	{
		Name = "Test",
		Value = 1,
	},
	{
		Name = "Test2",
		Value = 2,
	}
}

MY_TABLE_WITH_IDS = {
	[1005] = {
		Name = "Test",
		Value = 1,
	},
	[1006] = {
		Name = "Test2",
		Value = 2,
	}
}

MY_TABLE_INT_ARRAY = {
	[1005] = {
		Values = { 1, 2, 3 }
	},
	[1006] = {
		Values = { }
	},
	[1007] = {
		-- Values = { } -- not present
	},
}

MY_TABLE_SLICES = {
	[1005] = {
		Name = "Test",
		Value = {
			[100] = {
				{ FieldA = "DataA", FieldB = "DataB" },
			},
			[101] = {
				{ FieldA = "DataC", FieldB = "DataD" },
			},
		}
	},
	[1006] = {
		Name = "Test2",
		Value = {}
	}
}

MY_TABLE_INDEXED = {
	[1005] = { "Value1", 10 },
	[1006] = { "Value2", 20 },
}

MY_TABLE_INT = {
	Key1 = 1,
	Key2 = 2,
}

MY_TABLE_INT_METATABLE = {
	Key1 = 1,
	Key2 = 2,
	__newindex = function()
		error("unknown state")
	end
}
