package producer

import (
	"context"
	"encoding/json"

	"github.com/codfrm/cago/pkg/broker"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUpdateTopic = "collection.update"
)

type CollectionUpdateMsg struct {
	ID    primitive.ObjectID              `json:"id"`
	Exist map[primitive.ObjectID]struct{} `json:"exist"`
}

func PublishCollectionUpdate(ctx context.Context, id primitive.ObjectID, exist map[primitive.ObjectID]struct{}) error {
	body, err := json.Marshal(&CollectionUpdateMsg{
		ID:    id,
		Exist: exist,
	})
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, CollectionUpdateTopic, &broker2.Message{
		Body: body,
	})
}

func SubscribeCollectionUpdate(ctx context.Context,
	handler func(ctx context.Context, id primitive.ObjectID, exist map[primitive.ObjectID]struct{}) error) error {
	_, err := broker.Default().Subscribe(ctx, CollectionUpdateTopic, func(ctx context.Context, ev broker2.Event) error {
		msg := &CollectionUpdateMsg{}
		err := json.Unmarshal(ev.Message().Body, msg)
		if err != nil {
			return err
		}
		return handler(ctx, msg.ID, msg.Exist)
	})
	return err
}
