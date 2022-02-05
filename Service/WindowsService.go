package Service

import (
	"github.com/judwhite/go-svc"
	"log"
	"sync"
)

type program struct {
	wg   sync.WaitGroup
	quit chan struct{}
}

//TODO: replace all the 'log' crap with an actual logger.

func (p *program) Init(env svc.Environment) error {
	log.Printf("is win service? %v\n", env.IsWindowsService())
	return nil
}

func (p *program) Start() error {
	p.quit = make(chan struct{})

	p.wg.Add(1)
	go func() {
		//Do stuff
	}()

	return nil
}

func (p *program) Stop() error {
	log.Println("Stopping...")
	close(p.quit)
	p.wg.Wait()
	log.Println("Stopped.")
	
	return nil
}
