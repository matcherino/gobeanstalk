//  This package implements a single host thread safe connection pool for gobeanstalk
package gobeanstalk

//  Simple on demand connection pool
type Pool struct {
	addr string
	pool chan *Conn
}

// Makes a new connection pool addr is formated connection string (e.g. "localhost:11300")
func NewPool(addr string, size int) (*Pool, error) {
	pool := make([]*Conn, 0, size)
	for i := 0; i < size; i++ {
		client, err := Dial(addr)
		if err != nil {
			for _, client = range pool {
				client.Quit()
			}
			return nil, err
		}
		if client != nil {
			pool = append(pool, client)
		}
	}
	p := Pool{
		addr: addr,
		pool: make(chan *Conn, len(pool)),
	}
	for i := range pool {
		p.pool <- pool[i]
	}
	return &p, nil
}

// Retrieves an available beanstalk client. If there are none available it will create a new beanstalk connection
func (p *Pool) Get(tube string) (*Conn, error) {
	select {
	case conn := <-p.pool:
		return conn, nil
	default:
		return Dial(p.addr)
	}
}

// Returns a beanstalk connection back to the pool.  If the pool is full it will quit and discard the connection
func (p *Pool) Release(conn *Conn) {
	select {
	case p.pool <- conn:
	default:
		conn.Quit()
	}
}

// Quits all connections in pool
func (p *Pool) Empty() {
	var conn *Conn
	for {
		select {
		case conn = <-p.pool:
			conn.Quit()
		default:
			return
		}
	}
	p.pool = nil
}
