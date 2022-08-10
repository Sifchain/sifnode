//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package types

// clp module event types

const (
	EventTypeProcessedRemovalQueue = "processed_removal_queue"
	EventTypeQueueRemovalRequest   = "queue_removal_request"
	EventTypeDequeueRemovalRequest = "dequeue_removal_request"
	EventTypeProcessRemovalError   = "process_removal_error"
)
