local tblVar = arg[1];
local inFile = arg[2];
local forceTable = false;

if arg[3] == 'true' then
	forceTable = true;
end

JSON = (loadfile "lua/libs/json.lua")();
dofile(inFile);

local d = JSON:encode(_G[tblVar], forceTable);

print(d);
os.exit(0);
