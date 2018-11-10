package surlane

type Pool struct {
	queue chan []byte
	poolSize int
	bufSize int
}

func NewPool(poolSize int, bufSize int) *Pool {
	return &Pool{
		make(chan []byte, poolSize),
		poolSize,
		bufSize,
	}
}

func (pool *Pool) Borrow() (buffer []byte) {
	select {
	case buffer = <-pool.queue:
	default:
		buffer = make([]byte, pool.bufSize)
	}
	return
}

func (pool *Pool) GetBack(buffer []byte) {
	select {
	case  pool.queue<- buffer:
	default:
	}
	return
}