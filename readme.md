# Record/Replay commands

A simple lib (with cli) to record the outputs of an command with the same timedelays.

## Usage lib

```go
recorder := recmd.NewRecorder()
record := recorder.Record("curl duckduckgo.com",nil)
record2 := recorder.Record("bash", os.Stdin)


// record can be stored as a json


reader := record.Reader();
data, _ := io.ReadAll(reader);
```

## [WIP] Usage cli
