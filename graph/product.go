package graph

import (
	"auction-back/models"
	"auction-back/ports"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

func IsProductOwner(DB ports.DB, viewer models.User, product models.Product) error {
	owner, err := DB.Product().GetOwner(product)
	if err != nil {
		return fmt.Errorf("get owner: %w", err)
	}

	if owner.ID != viewer.ID {
		return ErrViewerNotOwner
	}

	return nil
}

func (r *productResolver) isProductOwnerOrManager(viewer models.User, obj models.Product) error {
	var errors error

	if err := IsProductOwner(r.DB, viewer, obj); err != nil {
		errors = multierror.Append(errors, err)
	} else {
		return nil
	}

	if err := r.RolePort.HasRole(models.RoleTypeManager, viewer); err != nil {
		errors = multierror.Append(errors, err)
	} else {
		return nil
	}

	return errors
}

// func (r *subscriptionResolver) ProductOffered(ctx context.Context) (<-chan *models.Product, error) {
// 	ch := make(chan *models.Product, 1)

// 	r.MarketLock.Lock()
// 	chan_id := randString(6)
// 	r.Market[chan_id] = ch
// 	r.MarketLock.Unlock()

// 	go func() {
// 		<-ctx.Done()
// 		r.MarketLock.Lock()
// 		delete(r.Market, chan_id)
// 		r.MarketLock.Unlock()
// 	}()

// 	return ch, nil
// }
