package parallel
//
//import (
//	"errors"
//	"sync"
//)
//
//const (
//	DataTypeList = iota
//	DataTypeListPointer
//	DataTypeChan
//)
//
//type Parallel struct {
//	Concurrency     int              // 并发数
//	Function        func(*Item)      // 执行的函数
//	DataType        int              // 处理的数据类型：DataTypeList、DataTypeListPointer、DataTypeChan
//	DataList        []interface{}    // 切片类型
//	DataListPointer *[]interface{}   // 切片指针类型
//	DataChan        chan interface{} // channel类型
//	queues          chan *Item       // 消费队列
//}
//
//func (p *Parallel) push(item *Item) {
//	p.queues <- item
//}
//
//func (p *Parallel) close() {
//	close(p.queues)
//}
//
///*
//并发运行
//入参：
// 1. 函数名
// 2. 并发数
// 3. 数据集
//   a. []interface{}类型：数据固定，不会变化
//   b. *[]interface{}类型：数据不固定，可能会变化
//   c. chan interface{}类型：数据不固定
//*/
//
//func (p *Parallel) checkData() error {
//	switch p.DataType {
//	case DataTypeList:
//		if p.DataList == nil {
//			return errors.New("DataList can not be nil")
//		}
//	case DataTypeListPointer:
//		if p.DataListPointer == nil {
//			return errors.New("DataListPointer can not be nil")
//		} else if *p.DataListPointer == nil {
//			return errors.New("DataListPointer's content can not be nil")
//		}
//	case DataTypeChan:
//		if p.DataChan == nil {
//			return errors.New("DataChan can not be nil")
//		}
//	default:
//		return errors.New("DataType is not valid")
//	}
//	return nil
//}
//
//func (p *Parallel) Run() error {
//	if err := p.checkData(); err != nil {
//		return err
//	}
//	if p.queues == nil {
//		p.queues = make(chan *Item)
//	}
//	// 数据异步压入队列
//	switch p.DataType {
//	case DataTypeList:
//		// 异步进行，数据压channel，压完后关闭channel
//		go func() {
//			for index, data := range p.DataList {
//				p.push(&Item{
//					Index:      index,
//					TotalCount: len(p.DataList),
//					Data:       data,
//				})
//			}
//			p.close()
//		}()
//	// 如果写的不及时，就存在导致认为已经读完了的问题，不推荐使用
//	case DataTypeListPointer:
//		// 异步进行，数据压channel，压完后关闭channel
//		go func() {
//			index := 0
//			for {
//				if index >= len(*p.DataListPointer) {
//					break
//				}
//				p.push(&Item{
//					Index:      index,
//					TotalCount: len(*p.DataListPointer),
//					Data:       (*p.DataListPointer)[index],
//				})
//			}
//			p.close()
//		}()
//	case DataTypeChan:
//		// 异步进行，数据压channel，压完后关闭channel
//		go func() {
//			index := 0
//			for data := range p.DataChan {
//				p.push(&Item{
//					Index:      index,
//					TotalCount: index,
//					Data:       data,
//				})
//			}
//			p.close()
//		}()
//	}
//
//	// 数据处理
//	var wg sync.WaitGroup
//	wg.Add(p.Concurrency)
//	for i := 0; i < p.Concurrency; i++ {
//		go func() {
//			for data := range p.queues {
//				p.Function(data)
//			}
//			wg.Done()
//		}()
//	}
//	wg.Wait()
//	return nil
//}
