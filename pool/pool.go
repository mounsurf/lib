package pool

import (
	"sync"
)

//https://github.com/dean2021/hackpool/blob/master/hackpool.go
type Pool struct {
	concurrency int
	queues      chan interface{}
	function    func(interface{})
}

func NewPool(concurrency int, function func(interface{})) *Pool {
	return &Pool{
		concurrency: concurrency,
		queues:      make(chan interface{}),
		function:    function,
	}
}
func (p *Pool) Push(data interface{}) {
	p.queues <- data
}

func (p *Pool) Close() {
	close(p.queues)
}

func (p *Pool) Run() {
	var wg sync.WaitGroup
	wg.Add(p.concurrency)
	for i := 0; i < p.concurrency; i++ {
		go func() {
			for v := range p.queues {
				p.function(v)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

//func out(x interface{}) {
//	fmt.Print(x," ")
//}
//
//func Dpool() {
//	p := NewPool(10, out)
//	go func() {
//		for e:=0;e<5;e++{
//			for i := 0; i < 10; i++ {
//				p.Push(i)
//			}
//			time.Sleep(time.Second * 1)
//			fmt.Println()
//		}
//		p.Close()
//	}()
//	p.Run()
//}
