package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/brharrelldev/BlockListAPI/blocklist"
	"github.com/brharrelldev/BlockListAPI/graph/generated"
	"github.com/brharrelldev/BlockListAPI/graph/model"
)

func (r *mutationResolver) Enqueue(ctx context.Context, ips []*string) ([]*model.Status, error) {
	bl, err := blocklist.NewBlockList(r.DB)
	if err != nil {
		return nil, fmt.Errorf("error instantiating new blocklist %v", err)
	}

	statuses, err := bl.Enqueue(ips)
	if err != nil {
		return nil, fmt.Errorf("error")
	}


	return statuses, err
}

func (r *queryResolver) GetIPDetails(ctx context.Context, ip string) (*model.IPAddress, error) {
	bl, err := blocklist.NewBlockList(r.DB)
	if err != nil {
		return nil, fmt.Errorf("error creating new blocklist instance %v", err)
	}

	ipDetal, err := bl.GetIPDetails(ip)
	if err != nil {
		return nil, fmt.Errorf("error getting ip detail %v", err)
	}

	return ipDetal, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
