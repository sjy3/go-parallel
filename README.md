# go-parallel
一个基于生产消费模式带超时时间的并发任务处理范式

go-parallel 是一个用于管理并发任务的 golang 库，基于这个工具，您可以将精力聚焦于开发具体业务逻辑，而不必担心并发任务的处理及超时。 go-parallel基于生产消费者模型实现并发任务的计算结果收集，支持并发流程的整体超时机制的设置。

例如：您需要并发的从多个接口拉取数据，并对数据进行汇聚。您可以创建一个并发任务，并添加多个function或者Process(如果需要参数传入的话)，并配置一个全局的超时时间。那么您在超时时间内会收到全部接口的结果，或者除超时之外的接口结果。

## 安装

```sh
go get github.com/sjy3/go-parallel
```

## 使用

对于大多数情况来说。通过 AppendFunc 来添加自己并发任务的逻辑，具体入下：

```golang

func main() {
	po := go_parallel.NewParallelObject()

	po.AppendFunc("", func() interface{} {
		time.Sleep(10 * time.Millisecond)
		return "func 1"
	})
	po.AppendFunc("", func() interface{} {
		time.Sleep(10 * time.Millisecond)
		return "func 2"
	})

	data, isTimeout := po.Run()
	var result []string
	for _,d := range data{
		result = append(result, d.Data.(string))
	}
	fmt.Println("isTimeout = ", isTimeout)
	fmt.Println(result)
	fmt.Println(len(result))
}

```

默认情况下，整个并发处理流程会等全部的子流程执行结束退出，isTimeout 恒为 false。此外，go-parallal 还支持设置超时时间，只需要在`po.Run()`之前设置参数, 如果存在超时子任务的话，isTimeout = true

```golang
	po.SetTimeout(5000 * time.Millisecond)
```

您也可以通过 context 上下文信息来处理超时，context 方式的好处是可以继承父任务的超时时间。

```golang
	ctx, cancel = context.WithTimeout(p.ctx, timeout)
    po.SetContext(ctx)
```

## 参数传入

当子任务需要依赖外部传入参数时，需要使用 `AppendProcess` 模式，具体来说，需要首先声明传入参数的执行对象，并实现 `ParallelProcessor` 接口，具体例子如下：

```golang
type process struct {
	param []int
}

func (p process) Do() interface{} {
	var result int
	for _, i := range p.param {
		result += i
	}
	return result
}

func main() {

	po := go_parallel.NewParallelObject()
	po.AppendProcess("", process{
		[]int{1,2,3,4,5},
	})
	po.AppendProcess("", process{
		[]int{5,6,7,8,9},
	})
	po.SetTimeout(5000 * time.Millisecond)

	data, isTimeout := po.Run()

	var result []int
	for _,d := range data{
		result = append(result, d.Data.(int))
	}

	fmt.Println("isTimeout = ", isTimeout)
	fmt.Println(result)
}
```

## 子任务区分

如果不同子任务的返回结果处理方式需要区分处理的情况，`AppendFunc` 和 `AppendProcess` 都支持子任务名称的设置（默认可以传空字符串），返回结果中会回传，可用于处理流程的区分。

```golang
func main() {
	po := go_parallel.NewParallelObject()
	po.AppendFunc("getName", func() interface{} {
		time.Sleep(10 * time.Millisecond)
		return "func 1"
	})
	po.AppendFunc("getId", func() interface{} {
		time.Sleep(10 * time.Millisecond)
		return 10
	})
	po.SetTimeout(5000 * time.Millisecond)
	data, isTimeout := po.Run()
	fmt.Println("isTimeout = ", isTimeout)

	for _,d := range data{
		switch d.Name{
		case "getName":
			fmt.Println("the result = ",d.Data.(string))
		case "getId":
			fmt.Println("the result = ",d.Data.(int))
		}
	}
}
```

## 注意事项

超时时间的设置，只能保证整体流程的正常退出，对于用户通过 `AppendFunc` 和 `AppendProcess` 声明的子任务，需要用户保证流程的正常退出，否则会有协程泄漏的风险存在。

