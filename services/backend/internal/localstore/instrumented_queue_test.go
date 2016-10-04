package localstore

import (
	"reflect"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"context"
)

func TestInstrumentQueue(t *testing.T) {
	Mocks = MockStores{}

	q := instrumentedQueue{}
	ctx := context.Background()

	// Enqueue
	want := &Job{Type: "test"}
	called := Mocks.Queue.MockEnqueue(t, want)
	if err := q.Enqueue(ctx, want); err != nil {
		t.Fatal(err)
	}
	if !*called {
		t.Error("Did not call underlying Enqueue")
	}

	// LockJob
	called, calledSuccess, calledError := Mocks.Queue.MockLockJob_Return(t, want)
	got, err := q.LockJob(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Job, want) {
		t.Errorf("unexpected LockedJob, got %+v, wanted %+v", got.Job, want)
	}
	if !*called {
		t.Error("Did not call underlying LockJob")
	}

	// LockJob.MarkSuccess
	*calledSuccess, *calledError = false, false
	err = got.MarkSuccess()
	if err != nil {
		t.Fatal(err)
	}
	if !*calledSuccess {
		t.Error("Did not call underlying LockJob.MarkSuccess")
	}
	if *calledError {
		t.Error("Called underlying LockJob.MarkError")
	}

	// LockJob.MarkError
	*calledSuccess, *calledError = false, false
	err = got.MarkError("test")
	if err != nil {
		t.Fatal(err)
	}
	if !*calledError {
		t.Error("Did not call underlying LockJob.MarkError")
	}
	if *calledSuccess {
		t.Error("Called underlying LockJob.MarkSuccess")
	}
}

func TestQueueStatsCollector(t *testing.T) {
	Mocks = MockStores{}

	stats := map[string]QueueStats{
		"a": QueueStats{
			NumJobs:          3,
			NumJobsWithError: 1,
		},
		"b": QueueStats{
			NumJobs: 1,
		},
	}
	Mocks.Queue.Stats = func(_ context.Context) (map[string]QueueStats, error) {
		return stats, nil
	}

	// We just check that we collect 4 stats, and don't actually check we
	// collect legit values.
	var (
		c     = newQueueStatsCollector(context.Background())
		ch    = make(chan prometheus.Metric)
		count = 0
	)
	go func() {
		c.Collect(ch)
		close(ch)
	}()
	for {
		select {
		case m := <-ch:
			if m == nil && count == 4 {
				return
			}
			if count > 4 || (m == nil && count != 4) {
				t.Fatalf("collected %d metrics, wanted 4", count)
			}
			count++
		case <-time.After(1 * time.Second):
			t.Fatal("expected collect timed out")
		}
	}
}
