package parallel

import (
	"errors"
	"reflect"
	"sync"
)

type parallel struct {
	concurrency int         // 并发数
	queues      chan *Item  // 消费队列
	function    func(*Item) // 执行的函数
}

func (p *parallel) push(item *Item) {
	p.queues <- item
}

func (p *parallel) close() {
	close(p.queues)
}

/*
并发运行
入参：
1. 函数名
2. 并发数
3. 数据集
  a. []interface{}类型：数据固定，不会变化
  b. *[]interface{}类型：数据不固定，可能会变化
  c. chan interface{}类型：数据不固定
*/
func Run(function func(*Item), concurrency int, dataSet interface{}) error {

	p := &parallel{
		function:    function,
		concurrency: concurrency,
		queues:      make(chan *Item),
	}
	// 数据异步压入队列
	switch reflect.TypeOf(dataSet).Kind() {
	case reflect.Slice:
		// 异步进行，数据压channel，压完后关闭channel
		go func() {
			// 参考 https://blog.csdn.net/sinat_35406909/article/details/104950795
			// 通过反射获取到数据，然后进行处理
			dataSlice := reflect.ValueOf(dataSet)
			for index := 0; index < dataSlice.Len(); index++ {
				p.push(&Item{
					Index:      index,
					TotalCount: dataSlice.Len(),
					Data:       dataSlice.Index(index).Interface(),
				})
			}
			p.close()
		}()
	case reflect.Chan:
		// 异步进行，数据压channel，压完后关闭channel
		go func() {
			index := 0
			v := reflect.ValueOf(dataSet)
			for {
				data, ok := v.Recv()
				if ok {
					p.push(&Item{
						Index:      index,
						TotalCount: index,
						Data:       data.Interface(),
					})
				} else {
					break
				}
				index++
			}
			p.close()
		}()
	// 如果写的不及时，就存在导致认为已经读完了的问题，不推荐使用
	case reflect.Ptr:
		// 异步进行，数据压channel，压完后关闭channel
		dataSlice := reflect.Indirect(reflect.ValueOf(dataSet))
		if dataSlice.Kind() == reflect.Slice {
			go func() {
				for index := 0; index < dataSlice.Len(); index++ {
					p.push(&Item{
						Index:      index,
						TotalCount: dataSlice.Len(),
						Data:       dataSlice.Index(index).Interface(),
					})

				}
				p.close()
			}()
		} else {
			p.close()
			return errors.New("dataSet: unsupported data type")
		}
	default:
		p.close()
		return errors.New("dataSet: unsupported data type")
	}

	// 数据处理
	var wg sync.WaitGroup
	wg.Add(p.concurrency)
	for i := 0; i < p.concurrency; i++ {
		go func() {
			for data := range p.queues {
				p.function(data)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}
