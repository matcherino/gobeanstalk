#gobeanstalk
Go [Beanstalkd](http://kr.github.io/beanstalkd/) client library.
Read the doc [here](http://godoc.org/github.com/iwanbk/gobeanstalk) .

## INSTALL
	go get github.com/iwanbk/gobeanstalk


## USAGE

### Producer
```go
import (
	"github.com/iwanbk/gobeanstalk"
	"log"
	"time"
)

func main() {
	conn, err := gobeanstalk.Dial("localhost:11300")
	if err != nil {
		log.Fatal(err)
	}

	id, err := conn.Put([]byte("hello"), 0, 10*time.Second, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Job id %d inserted\n", id)
}

```

### Consumer
```go
import (
	"github.com/iwanbk/gobeanstalk"
	"log"
)

func main() {
	conn, err := gobeanstalk.Dial("localhost:11300")
	if err != nil {
		log.Fatal(err)
	}
	for {
		j, err := conn.Reserve()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("id:%d, body:%s\n", j.ID, string(j.Body))
		err = conn.Delete(j.ID)
		if err != nil {
			log.Fatal(err)
		}
	}
}
```


### Connection pool
```go
import (
	"github.com/iwanbk/gobeanstalk"
	"log"
)

func main(){
	// Start a new connection pool (connection string, pool size)
	p, err := gobeanstalk.NewPool("localhost:11300", 10)
	if err != nil {
		log.Fatal(err)
	}

	// Retrive a pool member
	conn, err := p.Get()
	if err != nil {
		log.Fatal(err)
	}
	
	// Use the connection as you would a normal connection, if you use multiple tubes
	// always call Use on your connection, as a reused connection does not reset the tubei.
	err = conn.Use(testtube)
	if err != nil {
		log.Fatal(err)
		conn = nil
	}else{

		// Release the connection back into the pool, keep in mind that this connection 
		// is still using the "default" tube we specified from Use command.
		// If a connection is bad, do not release it to the pool, simply derefrence it.
		// The pool will aquire a healthy connection in its place on-demand.
		p.Release(conn)
	}

	// Empty the pool, this will close all pooled connections and destroy the pool
	p.Empty()
}
```


## Implemented Commands

Producer commands:

* use
* put

Worker commands:

* watch
* ignore
* reserve
* delete
* touch
* release
* bury

Other commands:

* stats-job
* quit


# Release Notes
Latest release is v0.3 that contains API changes, see release notes [here](https://github.com/iwanbk/gobeanstalk/blob/master/ReleaseNotes.txt)

## Author

* [Iwan Budi Kusnanto](http://iwan.my.id)
