# Basic usage
```go
c := rpc.New("http://localhost:8888")
result, err := c.Call("method", []interface{}{
	"param1",
	true,
	1,
})

```
