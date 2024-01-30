package consumer

import (
	"context"

	"github.com/dsp2b/dsp2b-go/internal/task/consumer/collection"

	"github.com/codfrm/cago/configs"
)

type Subscribe interface {
	Subscribe(ctx context.Context) error
}

// Consumer 消费者
func Consumer(ctx context.Context, cfg *configs.Config) error {
	subscribers := []Subscribe{
		&collection.Collection{},
	}
	for _, v := range subscribers {
		if err := v.Subscribe(ctx); err != nil {
			return err
		}
	}
	return nil
}
