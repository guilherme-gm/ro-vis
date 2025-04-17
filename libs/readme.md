This folder includes libraries we need to run ro-vis but has some limitation in being installed in the OS.

# Lua 5.1 32-bits

RO LUBs are compiled using Lua 5.1 32-bits, and we MUST use the same version in order to load LUB files
directly, without extracting them.

Ubuntu's apt-get no longer provides Lua 5.1 32-bits files, so you can't apt-get them.

My solution, for now, is to download the library and static link it to the ro-vis extractor.

Downloaded from: https://luabinaries.sourceforge.net/ (5.1.4 release 2)

> Linux26g4: Ubuntu 10.04 (x86) / Kernel 2.6 / gcc 4.4 / GTK 2.20
>
>

> [!NOTE]
> `liblua5.1.so` is renamed to `liblua5.1.so.bak` to avoid Go from trying to use it as
> a shared object, which would not work
