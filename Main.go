package main

import (
	"os"
	antidote "github.com/AntidoteDB/antidote-go-client"
	"fmt"
	"sync/atomic"
	"time"
	"math/rand"
	"strconv"
)

func main() { os.Exit(mainReturnWithCode()) }

var clientA int32
var clientB int32

func mainReturnWithCode() int {

	var bucketName []byte

	if len(os.Args) > 1 {
		bucketName = make([]byte, 8)
		for k,v := range os.Args[1:] {
			i, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("Error reading bucket name.")
				return 1
			}
			bucketName[k] = byte(i)
		}
	} else {
		bucketName = make([]byte, 8)
		rand.Seed(time.Now().UnixNano())
		rand.Read(bucketName)
	}

	bucket := antidote.Bucket{Bucket: bucketName}

	fmt.Printf("Start on bucket: %v\n", bucketName)

	if len(os.Args) == 1 {
		go incrementer(8101, &clientA, bucket)
		go incrementer(8201, &clientB, bucket)
	}

	go reader(8101, "A", bucket)
	reader(8201, "B", bucket)

	return 0
}




func incrementer(port int, counter *int32, bucket antidote.Bucket) error{
	host := antidote.Host{Name: "localhost", Port: port}
	client, err := antidote.NewClient(host)

	if err != nil {
		fmt.Println("Error creating Client.")
		return err
	}

	for {
		if err := incrementCounter(client, bucket); err != nil {
			fmt.Printf("Error incrementing counter %s \n", err)
		} else {
			atomic.AddInt32(counter, 1)
		}
	}
}

func reader(port int, reader string, bucket antidote.Bucket) error {
	host := antidote.Host{Name: "localhost", Port: port}
	client, err := antidote.NewClient(host)

	if err != nil {
		fmt.Println("Error creating Client.")
		return err
	}

	for {
		time.Sleep(10 * time.Second)
		if val,err := readCounter(client, bucket); err != nil {
			fmt.Printf("Error reading counter %s \n", err)
		} else {
			counterA := atomic.LoadInt32(&clientA)
			counterB := atomic.LoadInt32(&clientB)
			if reader == "A" {
				fmt.Printf("Readed value: %d on host %s, Total: %d, CounterA: %d, CounterB: %d, Received: %d, Pending: %d\n", val, reader, counterA+counterB, counterA, counterB, val-counterA, counterA+counterB-val)
			} else {
				fmt.Printf("Readed value: %d on host %s, Total: %d, CounterA: %d, CounterB: %d, Received: %d, Pending: %d\n", val, reader, counterA+counterB, counterA, counterB, val-counterB, counterA+counterB-val)
			}
		}
	}

}


func incrementCounter(client *antidote.Client, bucket antidote.Bucket) error {
	tx := client.CreateStaticTransaction()
	updateOp := antidote.CounterInc([]byte("key"), 1)
	return bucket.Update(tx, updateOp)
}


func readCounter(client *antidote.Client, bucket antidote.Bucket) (int32, error,){
	tx := client.CreateStaticTransaction()
	return bucket.ReadCounter(tx, []byte("key"))
}