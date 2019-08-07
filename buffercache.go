package buffercache

import (
"fmt"
"time"
)
 
type buf struct{
	data []byte
    time time.Time
}

//Buffercache Queue
type queue struct {
	buffers []*buf
	head int
	tail int
	count int
	memSize int
	bufSize int
	ageLimit int64
	maxMemSize int
	deletedBufCount int
}

var Queue *queue


//Function to initialize a buffer

func BufCacheInit(bsize int, qsize int , alimit int64) {
	if bsize > 2000 {
		fmt.Println("Warning: Too long buffers would lessen",
		" the capability to hold more buffers ")
	}
	if qsize > 2 * 1024 {
		fmt.Println("ERROR: Insufficent Memory resources on the system.",
		" See that it's less than 25% of your system memory")
		return 
	}
	Queue = &queue {buffers: make([]*buf,0), head: 0, tail: 0, 
	               count: 0, memSize: 0, bufSize: bsize, 
				   ageLimit: alimit*24*3600, maxMemSize: qsize, 
				   deletedBufCount:0}

	}

//Function to insert a buffer into the buffercache

func PutBuffer(b []byte) {
	if b == nil {
       fmt.Println("Nil buffer")
		return 
	}
    Queue.enqueue(b)
}
//Function to retrieve a buffer from the cache

func GetBuffer()[]byte {
	return Queue.dequeue()
}

//Displays the various statistics related to a buffer

func Stats() {
	fmt.Println("cache memory used = ", Queue.memSize) 
	fmt.Println("cache memory max size = ", Queue.maxMemSize) 
	fmt.Println("cache memory buffers holding time = ", Queue.ageLimit) 
	fmt.Println("cache memory buffer count = ", len(Queue.buffers)) 
	fmt.Println("cache memory deleted buffer count = ", Queue.deletedBufCount) 
}

//Modification of the buffer size and its age

func  ModifyBufCacheParams(msize int,age int64) {
	Queue.maxMemSize = msize
	Queue.ageLimit = age
	Queue.ageOut()
}

//Enqueuing a buffer into the queue

func (q *queue) enqueue (b  []byte) {
	if b == nil {
		return 
	}
	//Check the buffer size limit
	if len(b) > q.bufSize {
		fmt.Println("Warning: Too long buffers would lessen the capability",
		" to hold more buffers. This buffer len is ",len(b))
	}

	//check memory availability
	//Checking for buffer size limits
	if q.memSize + len(b) >= q.maxMemSize {
			fmt.Println("Buffer cache has reached 100%, old buffer" ,
			" will forcibly  aged out")
			//Create room for these many bytes on the queue
			q.dequeueEnough(len(b))
	} else if q.memSize + len(b) >= (q.maxMemSize*8)/10 {
			fmt.Println("Buffer cache has reached 80% threshold, after this" ,
			"old buffers may age out")
		}
	q.ageOut()
	buffer := &buf{data: b, time: time.Now()}
	q.buffers = append(q.buffers,buffer)
	q.memSize = q.memSize + len(b)
}

//Dequeuing a buffer from the queue

func (q *queue) dequeue () []byte {
	if len(q.buffers) == 0 {
		return nil
	}
		element := q.buffers[0]
		q.buffers = q.buffers[1:]
		q.memSize = q.memSize - len(element.data)
		return element.data
}

//Verifies the age restrictions of the buffers

func (q *queue) ageOut() {
	if len(q.buffers) == 0  {
		return
	}
	for _, v :=  range q.buffers {
		duration := time.Now().Sub(v.time)
		d := int64(duration)
	
        if d/1000 >= q.ageLimit {
			q.dequeue()
		    q.memSize = q.memSize - len(v.data)
			q.deletedBufCount++
		} else {
			break
		}
	}
}

func (q *queue) dequeueEnough(size int) {
	if len(q.buffers) == 0 {
		return
	}
	for element := q.dequeue(); size > 0; { 
		size = size - len(element) 
	}
}



