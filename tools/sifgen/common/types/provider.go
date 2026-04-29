
package types

import "context"

type Node struct {
	ID   string
	Name string
	IP   string
	// Add other relevant node properties
}

// Provider defines the interface for cloud provider interactions.
type Provider interface {
	GetNodes(ctx context.Context) ([]Node, error)
	CreateNode(ctx context.Context, name string) (*Node, error)
	DeleteNode(ctx context.Context, id string) error
}
