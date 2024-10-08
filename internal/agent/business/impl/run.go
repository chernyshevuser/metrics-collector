package impl

import (
	"context"
	"time"
)

func (s *svc) Run(ctx context.Context) {
	updateTicker := time.NewTicker(time.Duration(s.updateInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(s.sendInterval) * time.Second)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.closeCh:
				return
			case <-ctx.Done():
				return
			case <-updateTicker.C:
				s.logger.Infow(
					"update metrics",
					"status", "start",
				)
				s.collectMetrics()
				s.logger.Infow(
					"update metrics",
					"status", "finished",
				)
			}
		}
	}()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.closeCh:
				return
			case <-ctx.Done():
				return
			case <-updateTicker.C:
				s.logger.Infow(
					"update extra metrics",
					"status", "start",
				)
				s.collectExtraMetrics()
				s.logger.Infow(
					"update extra metrics",
					"status", "finished",
				)
			}
		}
	}()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.closeCh:
				return
			case <-ctx.Done():
				return
			case <-sendTicker.C:
				s.logger.Infow(
					"send metrics",
					"status", "start",
				)
				s.wg.Add(1)
				go func() {
					defer s.wg.Done()

					s.semaphore.Acquire()
					defer s.semaphore.Release()

					s.sendMetrics(ctx)
					s.logger.Infow(
						"send metrics",
						"status", "finished",
					)
				}()
			}
		}
	}()

	s.wg.Wait()
}
