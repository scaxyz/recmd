# Record/Replay command

A simple lib (with cli) to record the outputs of an command

## Usage lib

```go
recorder := recmd.NewRecorder()
record = recorder.Record("curl duckduckgo.com",nil)




reader := record.Reader();
data, _ := io.ReadAll(reader);



```