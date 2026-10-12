package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Capture compositor frames during real mouse/scroll/typing animations, rather
// than polling screenshots and repeating a few frames as a nominal 30 fps.
func startCapture(browser context.Context, output string, started time.Time, record *recording) (func() error, func() error, error) {
	type incoming struct {
		event   *page.EventScreencastFrame
		seconds float64
	}
	frames := make(chan incoming, 64)
	done := make(chan struct{})
	var failed atomic.Bool
	var stopped atomic.Bool
	var worker sync.WaitGroup
	chromedp.ListenTarget(browser, func(event any) {
		if frame, ok := event.(*page.EventScreencastFrame); ok && !stopped.Load() {
			select {
			case frames <- incoming{frame, time.Since(started).Seconds()}:
			default:
				failed.Store(true)
			}
		}
	})
	worker.Add(1)
	go func() {
		defer worker.Done()
		for {
			select {
			case <-done:
				return
			case frameEvent := <-frames:
				data, err := base64.StdEncoding.DecodeString(frameEvent.event.Data)
				name := fmt.Sprintf("frames/%06d.jpg", len(record.Frames))
				if err == nil {
					err = os.WriteFile(filepath.Join(output, name), data, 0600)
				}
				if err == nil {
					record.Frames = append(record.Frames, frame{File: name, Seconds: frameEvent.seconds})
				}
				if err != nil {
					failed.Store(true)
				}
				if err := chromedp.Run(browser, chromedp.ActionFunc(func(ctx context.Context) error {
					return page.ScreencastFrameAck(frameEvent.event.SessionID).Do(ctx)
				})); err != nil && !stopped.Load() {
					failed.Store(true)
				}
			}
		}
	}()
	if err := chromedp.Run(browser, chromedp.ActionFunc(func(ctx context.Context) error {
		return page.StartScreencast().WithFormat(page.ScreencastFormatJpeg).WithQuality(90).WithMaxWidth(1920).WithMaxHeight(1080).WithEveryNthFrame(1).Do(ctx)
	})); err != nil {
		close(done)
		worker.Wait()
		return nil, nil, err
	}
	check := func() error {
		if failed.Load() {
			return fmt.Errorf("compositor recording failed")
		}
		return nil
	}
	finish := func() error {
		stopped.Store(true)
		err := chromedp.Run(browser, chromedp.ActionFunc(func(ctx context.Context) error { return page.StopScreencast().Do(ctx) }))
		close(done)
		worker.Wait()
		record.Duration = time.Since(started).Seconds()
		if failed.Load() {
			return fmt.Errorf("compositor recording failed")
		}
		return err
	}
	return check, finish, nil
}
