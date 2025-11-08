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
