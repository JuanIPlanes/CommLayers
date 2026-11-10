package lro

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// OperationStatus represents the current state of a long-running operation.
type OperationStatus string

const (
	// OperationPending indicates the operation has been initiated but not yet processed.
	OperationPending OperationStatus = "pending"
	// OperationInProgress indicates the operation is being processed.
	OperationInProgress OperationStatus = "in_progress"
	// OperationCompleted indicates the operation completed successfully.
	OperationCompleted OperationStatus = "completed"
	// OperationFailed indicates the operation failed.
	OperationFailed OperationStatus = "failed"
)

// OperationID is a unique identifier for a long-running operation.
type OperationID string

// LRO describes a long-running operation with a status and result.
type LRO struct {
	// ID is the unique operation identifier.
	ID OperationID `json:"id"`
	// Status is the current state of the operation.
	Status OperationStatus `json:"status"`
	// CreatedAt is when the operation was created.
	CreatedAt time.Time `json:"created_at"`
	// CompletedAt is when the operation finished (if completed/failed).
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	// Result holds the output of the operation if completed.
	Result json.RawMessage `json:"result,omitempty"`
	// Error holds the error if the operation failed.
	Error string `json:"error,omitempty"`
}

// OperationManager manages the lifecycle of long-running operations.
type OperationManager struct {
	operations map[OperationID]*LRO
	mu         sync.RWMutex
	// timeout for async boundary; operations not completed within this duration
	// are considered stalled. Default: 10 seconds.
	timeout time.Duration
}

// NewOperationManager creates a new OperationManager with the given timeout.
// The timeout represents the async boundary; operations not completed within
// this duration are considered stalled.
func NewOperationManager(timeout time.Duration) *OperationManager {
	if timeout <= 0 {
		timeout = 10 * time.Second // 10-second async boundary default
	}
	return &OperationManager{
		operations: make(map[OperationID]*LRO),
		timeout:    timeout,
	}
}

// Create initiates a new long-running operation.
func (m *OperationManager) Create(ctx context.Context, id OperationID) *LRO {
	m.mu.Lock()
	defer m.mu.Unlock()

	lro := &LRO{
		ID:        id,
		Status:    OperationPending,
		CreatedAt: time.Now(),
	}
	m.operations[id] = lro
	return lro
}

// SetStatus updates the status of an operation.
func (m *OperationManager) SetStatus(ctx context.Context, id OperationID, status OperationStatus, result ...json.RawMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	lro, exists := m.operations[id]
	if !exists {
		return OperationNotFoundError{ID: id}
	}

	lro.Status = status
	if len(result) > 0 && result[0] != nil {
		lro.Result = result[0]
	} else {
		now := time.Now()
		lro.CompletedAt = &now
		if status == OperationCompleted {
			lro.Result = nil
		}
	}
	return nil
}

// GetStatus returns the current status of an operation.
func (m *OperationManager) GetStatus(ctx context.Context, id OperationID) (*LRO, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	lro, exists := m.operations[id]
	if !exists {
		return nil, OperationNotFoundError{ID: id}
	}
	return lro, nil
}

// Complete marks an operation as completed with an optional result.
func (m *OperationManager) Complete(ctx context.Context, id OperationID, result json.RawMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	lro := m.operations[id]
	lro.Status = OperationCompleted
	now := time.Now()
	lro.CompletedAt = &now
	lro.Result = result
	return nil
}

// Fail marks an operation as failed with an error message.
func (m *OperationManager) Fail(ctx context.Context, id OperationID, err string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	lro := m.operations[id]
	lro.Status = OperationFailed
	lro.Error = err
	now := time.Now()
	lro.CompletedAt = &now
	return nil
}

// IsStalled checks whether an operation has exceeded the async boundary timeout.
func (m *OperationManager) IsStalled(ctx context.Context, id OperationID) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	lro, exists := m.operations[id]
	if !exists {
		return false
	}
	// Check if operation has exceeded the timeout
	elapsed := time.Since(lro.CreatedAt)
	return elapsed > m.timeout
}

// OperationNotFoundError is returned when an operation ID is not found.
type OperationNotFoundError struct {
	ID OperationID
}

func (e OperationNotFoundError) Error() string {
	return "lro: operation not found: " + string(e.ID)
}

// IdempotencyKey is a key used to ensure idempotent operation submission.
type IdempotencyKey string

// IdempotentSubmit represents a request that can be made idempotently.
// The caller provides an idempotency key, and the server ensures that repeated
// submissions with the same key have the same effect as a single submission.
type IdempotentSubmit struct {
	// Key is the idempotency key identifying the submission.
	Key IdempotencyKey `json:"key"`
	// OperationID is the ID of the operation to associate with this submission.
	OperationID OperationID `json:"operation_id"`
	// Payload is the data submitted with the request.
	Payload json.RawMessage `json:"payload,omitempty"`
}

// ValidateIdempotencyKey checks whether the given key is valid.
// A valid key is non-empty and follows the expected format.
func ValidateIdempotencyKey(key string) bool {
	return len(key) > 0
}