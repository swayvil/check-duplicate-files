# check-duplicate-files
List duplicate files based on their hash signature. Optionnaly use "remove" parameter to <ins>move</ins> the duplicates to a folder named "REMOVED".

## Usage
```
go build
./check-duplicate-files <root directory path> [remove]
```

Only duplicate files are output:<br/>
```
<CRC32 Hash>;<Index>;<File path>
```
