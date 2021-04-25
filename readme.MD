# practice 1， errors are value.
go 中 error 的定义
```go
# http://golang.org/pkg/builtin/#error

type error interface {
	Error() string
}
```
基本的应用
```go
# http://golang.org/src/pkg/errors/errors.go
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
```
注意 New 方法返回的是地址，这在判断 哨兵Error 是否相等的时候非常必要，
因为 go 判断 struct 相等的方法是逐个地段判断相等。如果都相等，就值相等了，
而 不同 error 如果定义的时候 字符串 给相同了，会导致判断相同。
返回地址会杜绝这种情况。


# Sentinel Error
```go
if err == ErrSomething{
    // 
}

```
类似的 io.EOF， syscall.ENOENT

缺点：

不灵活，调用方要要依赖预先声明的error值，自己定义这种error 也会暴露更大的接口面积，成为API公共部分，在两个包创建了依赖；

可被破坏，
fmt.Errorf 携带一些上下文信息时， 会破坏调用者的 == 判断；调用者将被迫查看 error.Error() 信息，查看是否与特定字符串匹配。
不能依赖 error.Error() 方法的输出。
结论： 尽可能避免使用 sentinel errors.

# Error Types 携带上下文信息
```go
type MyError struct{
    Msg string
    File string
    Line int
}
func (e *MyError) Error string{
    return fmt.Sprintf("%s:%d:%s", e.File,e.Line, e.Msg)
}

func main() {
    err = MyError{}
    switch err := err.(type){
        case nil:
            // 
        case *MyError:
            fmt.PrintLn("xxx")
        default:
        // 
    }
}
```
例子有， os.PathError

虽然解决了一些痛点，带上了一些错误信息， 但是仍然不灵活，有sentinal error 耦合的缺点。
结论： 也尽量避免使用，至少避免将它变成公共api 的一部分。 本质仍然是 Sentinel error.

# Opaque errors， 不提供错误细节，只返回错误
只返回错误与否，不关注具体错误，但如何知道具体上下文呢，

Assert errors for behaviour, not type.

二分的判定错误与否，往往不够，有时候需要调查错误的性质。已确定操作是否合理。这时候可以断言错误实现了特定的行为，而不是断言错误的类型或值。

```go
package net
type Error interface{
    error

    Timeout() bool
    Temporary() bool
    ...
    ...
    ...
}
//使用

if nerr, ok:= err.(net.Error); ok && nerr.Temporary(){
    //
}

```
这么做，可以只暴露一个 Error, 然后判断 这个Error 的事件属性明确是哪种错误。而不用暴露很多 Error Type.
这就是 断言行为，而不是类型。

# 错误处理的tips

缩进的内容处理error.

暂存 error，集中处理。
```go
type errWriter struct{
    error
    io.Writer
}
func (e * errWriter) Write(buf []byte)(int, error){
    if e.err !=nil{
        return 0, e.err
    }
    var n int
    n, e.err = e.Writer.Writer(buf)
    return n, nil
}

使用起来：

func WriterResponse(w io.Writer, st Status, headers []Header, body io.Reader) error {
    ew := &errWriter{Writer:w}
    fmt.Fprintf(ew, "HTTP/1.1 %d %s\r\n", st.Code, st.Reason)
    for _, h := range headers{
        fmt.Fprintf(ew, "%s: %s\r\n", h.Key, h.Value)
        fmt.Fprint(ew, "\r\n")
        io.Copy(ew, body)
        return ew.err
    }
}
error 处理的代码复用了
```

处理错误（记录日志，或降级），或返回错误， 二者只能二选一

# 携带上下文信息

Wrap error
```go
pkg/errors
// 包装，携带信息,携带堆栈信息， %+v
errors.Wrap(err, "some msg")
// 获取原始错误
errors.Cause(err)
func Writer(){
_, err := w.Writer(buf)
// 如果没有错误， 这里返回 nil
return errors.Wrap(err, "writer failed)  
}

这个要注意， 错误不能重复 wrap 会携带多次堆栈信息， 所以这个库只适合应用级的开发包使用，不适合第三方包使用，第三方包使用了的话，用户无法知道错误被wrap了没有，再次warp 会有问题，所以不应该使用wrap， 应该只返回根因 error.
```

go 1.13 标准库 error 增加了 is  as
需要 error 类型 指定一个  Unwrap 方法
```go
func(e xER) Unwrap(){
    return e.error
}
```


