package tool

type Mutex struct {
	c chan int
}

func NewMutex() *Mutex {
	return &Mutex{make(chan int, 1)}
}

func (mu Mutex) Lock() { mu.c <- 1 }

func (mu Mutex) Unlock() { <-mu.c }

func (mu Mutex) Exec(fn func()) {
	mu.Lock()
	fn()
	mu.Unlock()
}
