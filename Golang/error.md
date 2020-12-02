**目录**
            
- [Error](#Error)         
    - [Error Type](##ErrorType)
        - [sentinel error](###sentinelError)
        - [Error types](###ErrorTypes)
        - [Opaque error](###OpaqueError)
    - [handing error](##handingError)
    - [Wrap error](##WrapError)
        - [wrap error使用场景](###wrapError使用场景)
    - [1.13后的error](##1.13后的error)
    - [Go 2版本Error处理](##Go2版本Error处理)
- [总结](#总结)


# Error
go的error只是一个interface。               
Go的处理异常逻辑是不引入exception，支持多参数返回。                         
如果一个函数返回了（value，error），不能对value进行假设，必须先判定error。如果对value不关心，那么也可以忽略error。            
panic只有在不可恢复的情况下，才会使用。              
诸如：索引越界，栈溢出等等。              

所以，error总结下来有以下几点。                  
- 简单
- 考虑失败，而不是成功            
- 没有隐藏的控制流
- 完全交给开发人员来控制error
- error are values

## ErrorType
### sentinelError 
预定义错误                                   
即，使用特定的值来表示error     

如 ErrSomething作为包级别的全局变量。            
```go
if err == ErrSomething{...}
```
类似的io.EOF,更底层的syscall.ENOENT。
使用sentinel值时最不灵活的错误处理策略，因为函数调用方必须使用==将结果与预先声明的值进行比较。当想要上下文的时候，就有问题了。因为要进行==判断，
返回的增加了上下文，判等操作就不好用了。                
就只能再被迫的去查看error.Error()方法的输出，来查看是否与特定的字符串匹配。                    

- 此类应该避免对err.Error()的依赖。                
- sentinel errors成为API的公共部分，对外暴露。
- 会在两个包里创建依赖。引用

结论：             
尽量少使用。              

### ErrorTypes
定义一个实现了error接口的自定义型，携带自己想要的信息。如：                    
```go
type MyError struct {
	Msg string
	File string
	Line int
}

func (e *MyError) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Msg)
}

func test() error {
	return &MyError{"something happened","xxx.go",42}
}
```
调用者可以用断言判断这个类型来使用。


### OpaqueError

**不透明error处理**
              
要求代码和调用者之间耦合最少。                 
就是你知道发生了错误，但不能看到错误内部。只需要返回错误，而不假设其内容。                   

调用者可以断言来判断是否是此类错误类型。如果是就调用其方法返回error。               
通常是定义一个对外非可见的interface。如                    
```go
type temporary interface {
	Temporary() bool
}
```

## handingError
- 判断不等于比判断等于要美观的多。 把异常的先抛出去。                         
- 消除不必要error判断。


## WrapError
没有上下文的错误，抛到最上层后，无法知道最终是哪报出来的错。需要把堆栈信息也携带者打印出来，就能定位到。                

你只应该处理error一次。避免同一个错误，一层层处处打日志。

- 错误要被日志记录。
- 应用程序处理错误，保证100%完整性。
- 之后不再报告当前错误。               


责任链模式：       
```go
    type withMessage struct {
    	cause error  //根因
    	msg   string
    }

    func Wrap(err error, message string) error {
    	if err == nil {
    		return nil
    	}
    	err = &withMessage{
    		cause: err,
    		msg:   message,
    	}
    	return &withStack{
    		err,
    		callers(),
    	}
    }
    
    func WithStack(err error) error {
    	if err == nil {
    		return nil
    	}
    	return &withStack{
    		err,
    		callers(),
    	}
    }
    
    type withStack struct {
    	error
    	*stack
    }
```
  
error 被封到withmessage中，withmessage被封到withstack中，最后使用wrap将报错return出去。

通过使用pkg/errors包，可以为错误值添加上下文。

### wrapError使用场景                
- 应用代码中，使用errors.New 或 errors.Errorf返回错误。                   
- 如果调用项目的其他包内的函数，通常简单直接返回。                 
- 如果和其他库协作，考虑使用errors.wrap， errors.wrapf保存堆栈信息。
- 直接返回错误，不要到处打日志。           
- 如果是在程序的顶部或者是在工作的goroutine顶部入口，使用%+v把堆栈详情记录。               
- 使用errors.cause获取root error，与sentinel error判定。                 

总结：
只有在application中使用，在标准库中只返回根因。                    
如果不打算处理这个error，就使用wrap往上抛。                                  
一旦这个错误被处理了，就不应该往上抛。                 

## 1.13后的error
自定义的error类型实现一个Unwrap方法，使根因能够得到返回。

检查错误的新函数                    
- Is
    errors.IS(err, ErrNotFound),会一层一层的去找根因。                 
- As
    errors.As(err &A) 也会一层层去找根因。                

可以还用%w向err添加附加信息。                   

## Go2版本Error处理
[Go 2 Error前瞻](https://go.googlesource.com/proposal/+/master/design/29934-error-values.md)


# 总结
错误的处理要基于自身业务。            
避免处处留日志的情况，Wrap error是非常好的一种方式。         
一定要将 error只处理一次牢记于心。            