# Go promise

Go version promise just for fun

## Usage

```go
NewPromise(func(resolve ResolveFunc, reject RejectFunc) {
    if true {
        resolve("ok")
    } else {
        reject(errors.New("failed"))
    }
}).Then(func(data interface{}) {
    fmt.Println(data)
}).Catch(func(err error) {
    fmt.Panic(err)
}).Do()
```
