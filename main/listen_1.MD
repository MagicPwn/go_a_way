# 本周作业

1. 我们在做数据库操作的时候， 比如DAO 层遇到一个sql.ErrorNoRows的时候，是否应该 Wrap 这个错误， 抛给上层，为什么，应该怎么做，请写出代码。

不应该 Wrap 这个错误， 应该直接返回

```go

if err == sql.ErrorNoRows:
    return err
```
