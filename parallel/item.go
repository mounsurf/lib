package parallel

type Item struct {
	Index      int         // 并发数
	TotalCount int         // 消费队列
	Data       interface{} // 执行的函数
}
