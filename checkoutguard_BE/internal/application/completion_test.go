package application

import (
	"context"
	"sync"
	"testing"
	"time"
)

type closingMarket struct {
	marketFixture
	entered, release chan struct{}
	once             sync.Once
}

func (b *closingMarket) Close() {
	b.once.Do(func() { close(b.entered) })
	<-b.release
}

type delayedRoot struct {
	noSpan
	entered, release chan struct{}
}

func (s delayedRoot) End() { close(s.entered); <-s.release }

type completionTrace struct {
	noTrace
	root  delayedRoot
	calls int
}

func (t *completionTrace) Workflow(ctx context.Context, _ string) (context.Context, Span) {
	t.calls++
	if t.calls == 1 {
		return ctx, t.root
	}
	return ctx, noSpan{}
}

func waitCompletionEvent(t *testing.T, event <-chan struct{}) {
	t.Helper()
	select {
	case <-event:
	case <-time.After(3 * time.Second):
		t.Fatal("completion event not reached")
	}
}

func TestTerminalRunIsPublishedAfterBrowserCloseAndCanStartDuringTraceExport(t *testing.T) {
	fixture := &closingMarket{marketFixture: marketFixture{state: "review"}, entered: make(chan struct{}), release: make(chan struct{})}
	trace := &completionTrace{root: delayedRoot{entered: make(chan struct{}), release: make(chan struct{})}}
	service := newService(newRepository(), testPlanner{next: func(context.Context, []Message) (Action, Usage, error) {
		return Action{Tool: "finish", Reason: "Review ready"}, Usage{}, nil
	}}, &marketGateway{state: MarketSessionState{Environment: "shopee", BrowserOpen: true}, browser: fixture})
	service.Tracer = trace
	t.Cleanup(func() {
		select {
		case <-fixture.release:
		default:
			close(fixture.release)
		}
		select {
		case <-trace.root.release:
		default:
			close(trace.root.release)
		}
		service.Close()
	})
	input := StartInput{Environment: "shopee", URL: "https://shopee.co.id/product/1/2", Goal: "Prepare checkout review"}
	first, err := service.Start(context.Background(), input, "")
	if err != nil {
		t.Fatal(err)
	}
	waitCompletionEvent(t, fixture.entered)
	progress, _ := service.Get(context.Background(), first.ID)
	if progress.Status != "running" || progress.Finished != nil || len(progress.Steps) != 2 {
		t.Fatal("terminal result published while browser was still closing")
	}
	close(fixture.release)
	waitCompletionEvent(t, trace.root.entered)
	final, _ := service.Get(context.Background(), first.ID)
	if final.Status != "ready_for_review" || final.Finished == nil {
		t.Fatal("final result missing after browser cleanup")
	}
	// The prior root remains blocked in export; it must no longer prevent a
	// fresh run or clear the fresh run's admission when export eventually ends.
	second, err := service.Start(context.Background(), input, "")
	close(trace.root.release)
	service.Close()
	if err != nil || second.ID == first.ID {
		t.Fatal("completed run prevented next run", err)
	}
}
