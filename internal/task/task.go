package task

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusRunning Status = "running"
	StatusDone    Status = "done"
	StatusFailed  Status = "failed"
)

type TaskDTO struct {
	ID        string
	Status    Status
	CreatedAt time.Time
	Duration  float64
	Result    *string
	Error     *string
}
type Task struct {
	ID         string    `json:"id"`
	Status     Status    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	Duration   float64   `json:"duration"` // в секундах
	Result     *string   `json:"result,omitempty"`
	Error      *string   `json:"error,omitempty"`
	cancelFunc context.CancelFunc
	mu         sync.Mutex
}

func NewTask() *Task {
	ctx, cancel := context.WithCancel(context.Background())
	t := &Task{
		ID:         uuid.NewString(),
		Status:     StatusPending,
		CreatedAt:  time.Now(),
		cancelFunc: cancel,
	}
	go t.run(ctx)
	return t
}

func (t *Task) run(ctx context.Context) {
	t.updateStatus(StatusRunning)
	start := time.Now()

	// Симуляция долгой работы (3-5 минут)
	select {
	case <-time.After(time.Duration(180+rand.Intn(120)) * time.Second):
		res := "ok"
		t.setResult(res, nil)
	case <-ctx.Done():
		err := "task cancelled"
		t.setResult("", &err)
	}
	t.setDuration(time.Since(start).Seconds())
}

func (t *Task) updateStatus(status Status) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Status = status
}

func (t *Task) setResult(result string, err *string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if err != nil {
		t.Status = StatusFailed
		t.Error = err
	} else {
		t.Status = StatusDone
		t.Result = &result
	}
}

func (t *Task) setDuration(d float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Duration = d
}

func (t *Task) Cancel() {
	t.cancelFunc()
}

func (t *Task) ToDTO() TaskDTO {
	t.mu.Lock()
	defer t.mu.Unlock()

	return TaskDTO{
		ID:        t.ID,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		Duration:  t.Duration,
		Result:    t.Result,
		Error:     t.Error,
	}
}
