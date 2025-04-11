package infrastructure

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"leonardovee.dev/command/internal/command"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

const version = 1
const topic = "arn:aws:sns:us-east-1:000000000000:command-events"

type Event struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Body      command.Command `json:"body"`
	CreatedAt string          `json:"created_at"`
	Version   int             `json:"version"`
}

type SNS struct {
	logger *slog.Logger
	client *sns.Client
}

func NewSNS(logger *slog.Logger, client *sns.Client) *SNS {
	return &SNS{
		logger: logger,
		client: client,
	}
}

func (s *SNS) Publish(ctx context.Context, command command.Command) {
	event := Event{
		Version:   version,
		ID:        command.GetId(),
		Name:      string(command.GetName()),
		Body:      command,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	message, err := json.Marshal(event)
	if err != nil {
		s.logger.Error("failed to marshal event", "name", command.GetName(), "error", err)
		return
	}

	publishInput := sns.PublishInput{
		TopicArn: aws.String(topic),
		Message:  aws.String(string(message)),
	}

	_, err = s.client.Publish(ctx, &publishInput)
	if err != nil {
		s.logger.Error("failed to publish command", "name", command.GetName(), "error", err)
		return
	}

	s.logger.Info("command published", "name", command.GetName())
}
