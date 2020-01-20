package main

import (
	"fmt"
	"github.com/icrowley/fake"
	"os"
	"sync"
	"time"
)

type Worker interface {
	Handle()
}

type TimeWorker struct {}

type FileReadWork struct {
	fileName string
	data string
}

type Exit struct {}

type RoutinesCloser struct {
	stopQ chan Worker
}

type TimeOutWorker struct {}

func (TimeWorker) Handle()  {
	fmt.Println(time.Now())
}

func (routinesCloser RoutinesCloser) Handle()  {
	fmt.Println("Close all routines")
	//close(routinesCloser.stopQ)
}

func (Exit) Handle() {
	fmt.Println("exitWorker is in progress")
	os.Exit(1)
}

func (frw FileReadWork) Handle()  {
	fmt.Println("Read and write file routine")

	filePath := "./files/" + frw.fileName + time.Now().Format("2006.01.02 15:04:05")

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	check(err)

	_, err = f.Write([]byte(frw.data))
	check(err)

	f.Close()
}

func (TimeOutWorker) Handle() {
	fmt.Println("timeOutWorker is in progress")
	TimeWorker{}.Handle()
	time.Sleep(time.Second)
	TimeWorker{}.Handle()

}

func worker(wg *sync.WaitGroup, taskQ chan Worker) {
	for task := range taskQ {
		task.Handle()
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
		//RoutinesCloser{workQ},
		Exit{},
		TimeOutWorker{},

	}
	fmt.Println(len(tasks))
	go worker(&wg, workQ)

	for _, task := range tasks {
		wg.Add(1)
		workQ <- task
	}

	wg.Wait()
	fmt.Println("Never reach this ")
}

func check (e error) {
	if e != nil {
		panic(e)
	}
}
