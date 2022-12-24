package utils

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/logger"
	"go.uber.org/zap"
)

func WithGracefulShutdown(cancel func(), jobs ...func()) {
	wg := sync.WaitGroup{}
	wg.Add(len(jobs))
	for _, job := range jobs {
		job := job
		go func() {
			defer wg.Done()
			job()
		}()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Ready for stop signal")
	sig := <-sigChan

	logger.Info("Received %v, gracefully shutting down...\n", zap.String("signal", sig.String()))
	cancel()

	WaitWithTimeout(&wg, 10*time.Second)
}

func WaitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) (timedOut bool) {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}
