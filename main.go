package main

import (
	"context"
	"fmt"
	"github.com/icrowley/fake"
	"os"
	"sync"
	"time"
)

var done chan bool

type Worker interface {
	Handle(ctx context.Context)
}

type TimeWorker struct {}

func (TimeWorker) Handle(ctx context.Context)  {
	fmt.Println("Current time file routine")
	fmt.Println(time.Now())
}

type FileReadWork struct {
	fileName string
	data string
}

func (frw FileReadWork) Handle(ctx context.Context)  {
	fmt.Println("Read and write file routine")

	filePath := "./files/" + frw.fileName + time.Now().Format("2006.01.02 15:04:05")

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	check(err)

	_, err = f.Write([]byte(frw.data))
	check(err)

	f.Close()
}

type Exit struct {}

func (Exit) Handle(ctx context.Context) {
	fmt.Println("exitWorker is in progress")
	os.Exit(1)
}

type RoutinesCloser struct {}

func (routinesCloser RoutinesCloser) Handle(ctx context.Context) {

	fmt.Println("RoutinesCloser is in progress")
	ctx, cancel := context.WithCancel(ctx)

	cancel()

	select {
	case <-ctx.Done():
		fmt.Println("Gracefully exit")
		fmt.Println(ctx.Err())
		return
	default:
	}

}

type TimeOutWorker struct {}

func (TimeOutWorker) Handle(ctx context.Context) {
	fmt.Println("Time Out Worker file routine")
	fmt.Println(time.Now())
	time.Sleep(time.Second)
	fmt.Println(time.Now())

}

func worker(wg *sync.WaitGroup, taskQ chan Worker) {

	ctx := context.Background()

	for task := range taskQ {
		task.Handle(ctx)
		wg.Done()
	}

}

func main()  {

	var wg sync.WaitGroup

	workQ := make(chan Worker)

	fileData := FileReadWork{
		fileName: fake.Word(),
		data:     fake.Sentence(),
	}

	tasks := []Worker {
		TimeWorker{},
		fileData,
		RoutinesCloser{},
		//Exit{},
		TimeOutWorker{},

	}

	fmt.Println(len(tasks))

	go worker(&wg, workQ)

	for _, task := range tasks {

		wg.Add(1)
		workQ <- task

	}

	wg.Wait()
	close(workQ)
}

func check (e error) {
	if e != nil {
		panic(e)
	}
}

