# 有没有遇到复杂的查询，数据量大了怎么办，以及提高查询效率、如何实现报警推送


# 场景 服务端大文件，http多线程下载



# error 处理
pkg/errors

# 遇到过的问题
如： context的使用,使用defer cancel()了，却在野生goroutine中进行了使用，导致报context canceled的错误。                 
例               
```go
    ctx, cancel := context.WithCancel(ctx.BackGround())
    defer cancel()
    
    group.Go(func(){} error){
    	
    }
    
    go func(){
        ctx直接使用 //错误的姿势  context canceled
        context.BackGround() //正确的姿势
    }()
    

```




