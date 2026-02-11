package launcher

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/ra341/glacier/pkg/syncmap"
	"github.com/rs/zerolog/log"
)

type Process struct {
	ctx    context.Context
	cancel context.CancelFunc
	cmd    *exec.Cmd
}

type Service struct {
	exeMap syncmap.Map[string, *Process]
}

func New() *Service {
	return &Service{}
}

func (s *Service) Launch(exe string) error {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, exe)
	err := cmd.Start()
	if err != nil {
		cancel()
		return err
	}

	pr := &Process{
		ctx:    ctx,
		cancel: cancel,
		cmd:    cmd,
	}

	go func() {
		s.supervisor(exe, pr)
	}()

	s.exeMap.Store(exe, pr)

	return nil
}

func (s *Service) supervisor(exe string, process *Process) {
	defer func() {
		process.cancel()
		s.exeMap.Delete(exe)
	}()

	err := process.cmd.Wait()
	if process.ctx.Err() != nil {
		log.Debug().Str("exe", exe).Msg("process closed via context")
		return
	}

	if err != nil {
		log.Warn().Err(err).Str("exe", exe).Msg("process exited")
		return
	}

	log.Info().Str("exe", exe).Msg("Process exited independently")
}

func (s *Service) Running(ctx context.Context, exe string) error {
	val, ok := s.exeMap.Load(exe)
	if !ok {
		return fmt.Errorf("exe is not running")
	}

	select {
	case <-ctx.Done():
		// the request context
		return nil
	case <-val.ctx.Done():
		// or wait for proces to be done
		return nil
	}
}
