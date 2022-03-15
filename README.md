# check-duplicate-files
List duplicate files based on their hash signature

## Usage
```
go build
./check-duplicate-files <root directory path>
```

Output only for duplicates:<br/>
```
<CRC32 Hash>;<Index>;<File path>
```
