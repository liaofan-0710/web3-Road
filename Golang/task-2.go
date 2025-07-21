package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 1. 题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，
// 然后在主函数中调用该函数并输出修改后的值。
//   - 考察点 ：指针的使用、值传递与引用传递的区别。

func testPointer1(num *int) int {
	*num += 10
	return *num
}

// 2. 题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
//   - 考察点 ：指针运算、切片操作。
func testPointer2(nums *[]int) []int {
	for i := 0; i < len(*nums); i++ {
		(*nums)[i] *= 2
	}
	return *nums
}

// 1. 题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
//   - 考察点 ： go 关键字的使用、协程的并发执行。

func testGoroutine1() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for count := 1; count <= 10; count++ {
			if count%2 == 0 {
				fmt.Print("偶数：", count, " ")
				count++
			}
		}
	}()
	go func() {
		defer wg.Done()
		for count := 1; count <= 10; count++ {
			if count%2 != 0 {
				fmt.Print("奇数：", count, " ")
				count++
			}
		}
		fmt.Println()
	}()
	wg.Wait()
}

// 2. 题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，
// 同时统计每个任务的执行时间。
//   - 考察点 ：协程原理、并发任务调度。

type Task func()

type TaskResult struct {
	ID        int
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

func testGoroutine2(tasks []Task) {
	results := make(chan TaskResult, len(tasks))
	wg := sync.WaitGroup{}
	wg.Add(len(tasks))
	for id, task := range tasks {
		go func(id int, task Task) {
			defer wg.Done()
			// 记录开始时间
			start := time.Now()

			// 执行任务
			task()

			// 记录结束时间并计算持续时间
			end := time.Now()
			duration := end.Sub(start)

			// 发送结果到通道
			results <- TaskResult{
				ID:        id,
				StartTime: start,
				EndTime:   end,
				Duration:  duration,
			}

		}(id, task)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Printf("%2d | %s | %s | %v\n",
			result.ID,
			result.StartTime.Format("2006-01-02 15:04:05.000"),
			result.EndTime.Format("2006-01-02 15:04:05.000"),
			result.Duration.Round(time.Millisecond),
		)
	}
}

// 1. 题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，
// 实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
//   - 考察点 ：接口的定义与实现、面向对象编程风格。

func (rectangle Rectangle) Area() float64 {
	return rectangle.height * rectangle.width
}

func (circle Circle) Perimeter() float64 {
	return circle.radius * 2 * 3.14159
}

type Rectangle struct {
	width, height float64
}

type Circle struct {
	radius float64
}

type Shape interface {
	Area() float64
	Perimeter() float64
}

// 2. 题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
//   - 考察点 ：组合的使用、方法接收者。

type Person struct {
	Name string
	Age  uint8
}

type Employee struct {
	Person
	EmployeeID int
}

func (e Employee) PrintInfo() Person {
	if e.EmployeeID == 1 {
		return e.Person
	} else if e.EmployeeID == 2 {
		return e.Person
	}
	return Person{}
}

// 题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，
// 另一个协程从通道中接收这些整数并打印出来。
// 考察点 ：通道的基本使用、协程间通信。
func testChan1() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	ch := make(chan int, 10)
	go func() {
		defer wg.Done()
		defer close(ch)
		for i := 1; i <= 10; i++ {
			ch <- i
		}
	}()
	go func() {
		defer wg.Done()
		for num := range ch {
			fmt.Print(num, " ")
		}
	}()
	wg.Wait()
}

func testChan2() {
	ch := make(chan int, 10)
	done := make(chan bool)
	go func() {
		defer close(ch)
		for i := 1; i <= 10; i++ {
			ch <- i
		}
	}()
	go func() {
		for num := range ch {
			fmt.Print(num, " ")
		}
		done <- true
	}()
	<-done
}

// 2. 题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
//   - 考察点 ：通道的缓冲机制。

func testChan3() {
	ch := make(chan int, 100)
	done := make(chan struct{})
	go func() {
		defer close(ch)
		for i := 1; i <= 100; i++ {
			ch <- i
		}
	}()
	go func() {
		for num := range ch {
			fmt.Print(num, " ")
		}
		close(done)
	}()
	<-done
}

// 1. 题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。
// 启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
//   - 考察点 ： sync.Mutex 的使用、并发数据安全。
func testLock1() {
	mx := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(10)
	count := 0
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for iota1 := 1; iota1 <= 1000; iota1++ {
				mx.Lock()
				count++
				mx.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Println("count:", count)
}

// 2. 题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。
// 启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
//   - 考察点 ：原子操作、并发数据安全。
func testLock2() {
	wg := &sync.WaitGroup{}
	wg.Add(10)
	count := int64(0)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for iota1 := 1; iota1 <= 1000; iota1++ {
				atomic.AddInt64(&count, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Println("count:", atomic.LoadInt64(&count))
}

func main() {
	// 题目 1.1
	num := 0
	fmt.Println(testPointer1(&num))
	// 题目 1.2
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println(testPointer2(&nums))
	// 题目 2.1
	testGoroutine1()
	// 题目 2.2
	tasks := []Task{
		func() { time.Sleep(300 * time.Millisecond) }, // 任务0
		func() { time.Sleep(200 * time.Millisecond) }, // 任务1
		func() { time.Sleep(100 * time.Millisecond) }, // 任务2
		func() { time.Sleep(400 * time.Millisecond) }, // 任务3
		func() { time.Sleep(150 * time.Millisecond) }, // 任务4
	}
	testGoroutine2(tasks)
	// 题目 3.1
	rect := Rectangle{width: 5, height: 3} // 长 宽
	fmt.Println(rect.Area())               // 返回总面积
	circle := Circle{radius: 3.5}          // 半径
	fmt.Println(circle.Perimeter())        // 返回周长
	// 题目 3.2
	emp := Employee{Person: Person{Name: "liaofan", Age: 25}, EmployeeID: 1}
	fmt.Println(emp.PrintInfo())
	// 题目 4.1
	testChan1()
	testChan2()
	// 题目 4.2
	testChan3()
	// 题目 5.1
	testLock1()
	// 题目 5.2
	testLock2()
}
