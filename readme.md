# Record/Replay command

A simple lib (with cli) to record the outputs of an command

## Usage lib

```go
recorder := recmd.NewRecorder()
record = recorder.Record("curl duckduckgo.com",nil)

// record can be stored as a json


reader := record.Reader();
data, _ := io.ReadAll(reader);
```

## [WIP] Usage cli