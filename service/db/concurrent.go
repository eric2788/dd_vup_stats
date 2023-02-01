package db

import "sync"

type (
	ParallelStream struct {
		errs  chan error
		stmts []Statement
		wg    *sync.WaitGroup
	}

	Statement func() error
)

func NewParallelStream() *ParallelStream {
	return &ParallelStream{
		wg: &sync.WaitGroup{},
	}
}

func (p *ParallelStream) AddStmt(stmt Statement) {
	p.stmts = append(p.stmts, stmt)
	p.wg.Add(1)
}

// Run blocking wait for all statements to finish
// and return the first error if any
func (p *ParallelStream) Run() error {
	p.errs = make(chan error, len(p.stmts))
	for _, stmt := range p.stmts {
		go func(stmt Statement) {
			defer p.wg.Done()
			p.errs <- stmt()
		}(stmt)
	}
	p.wg.Wait()
	close(p.errs)
	for err := range p.errs {
		if err != nil {
			return err
		}
	}
	return nil
}
