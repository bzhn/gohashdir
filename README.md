## Usage

```
.\gohashdir.exe C:\Users\username\Downloads\
```

The output will looks like following:

```
Path:		C:\Users\username\Downloads\
Name:		linuxmint-21-xfce-64bit.iso
xxhash64:	a9836e2d4ace9c44
sha256:		3ad001dc15cb661c6652ce1d20ecdc85a939fa0b4b9325af5d0c65379cc3b17e
Size:		2415585280
Time:		6304821800
--------
Path:		C:\Users\username\Downloads\
Name:		ubuntu-22.04.1-live-server-amd64.iso
xxhash64:	78e2c3c04f773030
sha256:		10f19c5b2b8d6db711582e0e27f5116296c34fe4b313ba45f9b201a5007056cb
Size:		1474873344
Time:		4814139300
--------
```

## TODO:
Flags to be added (there's no such flags yet):
-r [<n>]    Scan files recursively. Optional the depth may be specified
-h          Human readable output of elapsed time and size of files
-S          Sort by size
--json      Make output in JSON

