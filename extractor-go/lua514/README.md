This folder contains Lua 5.1.4 source with small modifications to MakeFile for ro-vis.

# Build

## Pre-requisites
If you are on a 32-bit system, you should be good to go.

If you are on 64-bit system, you need to install 32-bit version of readline:

```bash
sudo apt-get -y install lib32readline-dev
```


## Build

On 32-bit systems: (Linux)
```bash
make linux test
```

On 64-bit systems: (Linux)
```bash
make linux32 test
```

For other OSes, try running `make all` to see the available platforms.

Note that you MUST have a 32-bits static library as output.

After the build is complete, ensure it is a 32-bits build:

```bash
objdump -f src/liblua.a | grep ^architecture
```

Should output:

```
architecture: i386, flags 0x00000011:
```

With that, you should be good to run the go program.


# Modifications

1. Added `linux32` to PLATS in `Makefile`
2. Added `linux32` to PLATS in `src/Makefile`
3. Added `linux32` build target in `src/Makefile`:

```makefile
linux32: # Custom platform to build 32-bits on 64-bits systems
	$(MAKE) all MYCFLAGS="-DLUA_USE_LINUX -m32" MYLDFLAGS="-m32" MYLIBS="-Wl,-E -ldl -lreadline -lhistory -lncurses"

```

The objective of this change is to provide the `-m32` flag to `MYCFLAGS` and `MYLDFLAGS`.
This enables a 32-bits build on 64-bits systems.
