package gcp

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"cloud.google.com/go/pubsub"
)

// MessageRepository defines the interface for interacting with GCP Pub/Sub.
type MessageRepository interface {
	Publish(ctx context.Context, event string, data []byte) error
	PublishEvent(ctx context.Context, event string, data []byte) error
}

type messageRepository struct {
	log   *slog.Logger
	topic *pubsub.Topic
}

func NewMessageRepository(ctx context.Context, log *slog.Logger, projectID string) (*messageRepository, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	topic := client.Topic("event-bus")

	// Disable batching
	topic.PublishSettings.CountThreshold = 1
	topic.PublishSettings.DelayThreshold = 0

	return &messageRepository{
		log:   log,
		topic: topic,
	}, nil
}

func (r *messageRepository) Publish(ctx context.Context, event string, data []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.topic.Publish(ctx, &pubsub.Message{
		Data: data,
		Attributes: map[string]string{
			"event": event,
		},
	})

	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	r.log.Debug(fmt.Sprintf("published: %s", id))
	return nil
}

// Publish sends a message to the specified topic with the event set as an attribute.
func (r *messageRepository) PublishEvent(ctx context.Context, event string, data []byte) error {
	return r.Publish(ctx, event, data)
}
