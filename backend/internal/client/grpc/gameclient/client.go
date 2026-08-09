package gameclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	avitoshav1 "github.com/guitaramust-sudo/Avitosha/app/backend/internal/gen/avitosha/v1"
	"github.com/guitaramust-sudo/Avitosha/app/backend/internal/model"
	"github.com/guitaramust-sudo/Avitosha/app/backend/internal/realtime"
	internalrpc "github.com/guitaramust-sudo/Avitosha/app/backend/internal/rpc"
	"github.com/guitaramust-sudo/Avitosha/app/backend/internal/usecase"
)

type Client struct{ rpc avitoshav1.GameServiceClient }

func New(rpc avitoshav1.GameServiceClient) *Client { return &Client{rpc: rpc} }

func (c *Client) EnsureProfile(ctx context.Context, userID uuid.UUID, at time.Time) (usecase.GameProfile, error) {
	response, err := c.rpc.EnsureProfile(ctx, userAt(userID, at))
	return decode[usecase.GameProfile](response, err)
}

func (c *Client) RenamePet(ctx context.Context, userID uuid.UUID, name string, at time.Time) (usecase.GameProfile, error) {
	response, err := c.rpc.RenamePet(ctx, &avitoshav1.RenamePetRequest{UserId: userID.String(), Name: name, At: formatTime(at)})
	return decode[usecase.GameProfile](response, err)
}

func (c *Client) ListTasks(ctx context.Context, userID uuid.UUID, at time.Time) ([]model.TaskProgress, error) {
	response, err := c.rpc.ListTasks(ctx, userAt(userID, at))
	return decode[[]model.TaskProgress](response, err)
}

func (c *Client) GetTask(ctx context.Context, userID, taskID uuid.UUID, at time.Time) (model.TaskProgress, error) {
	response, err := c.rpc.GetTask(ctx, taskRequest(userID, taskID, at))
	return decode[model.TaskProgress](response, err)
}

func (c *Client) GetTaskAdvice(ctx context.Context, userID, taskID uuid.UUID, at time.Time) (usecase.TaskAdvice, error) {
	response, err := c.rpc.GetTaskAdvice(ctx, taskRequest(userID, taskID, at))
	return decode[usecase.TaskAdvice](response, err)
}

func (c *Client) GetRoom(ctx context.Context, userID uuid.UUID, at time.Time) ([]model.RoomItemProgress, error) {
	response, err := c.rpc.GetRoom(ctx, userAt(userID, at))
	return decode[[]model.RoomItemProgress](response, err)
}

func (c *Client) GetStory(ctx context.Context, userID uuid.UUID, at time.Time) (model.StorySnapshot, error) {
	response, err := c.rpc.GetStory(ctx, userAt(userID, at))
	return decode[model.StorySnapshot](response, err)
}

func (c *Client) GetDailySummary(ctx context.Context, userID uuid.UUID, at time.Time) (usecase.DailySummary, error) {
	response, err := c.rpc.GetDailySummary(ctx, userAt(userID, at))
	return decode[usecase.DailySummary](response, err)
}

func (c *Client) GetLeaderboard(ctx context.Context, userID uuid.UUID, limit int, at time.Time) (usecase.Leaderboard, error) {
	response, err := c.rpc.GetLeaderboard(ctx, &avitoshav1.LeaderboardRequest{UserId: userID.String(), Limit: int32(limit), At: formatTime(at)})
	return decode[usecase.Leaderboard](response, err)
}

func (c *Client) GetAchievements(ctx context.Context, userID uuid.UUID, at time.Time) ([]model.AchievementProgress, error) {
	response, err := c.rpc.GetAchievements(ctx, userAt(userID, at))
	return decode[[]model.AchievementProgress](response, err)
}

func (c *Client) GetRewardBalances(ctx context.Context, userID uuid.UUID, at time.Time) ([]model.RewardBalance, error) {
	response, err := c.rpc.GetRewardBalances(ctx, userAt(userID, at))
	return decode[[]model.RewardBalance](response, err)
}

func (c *Client) GetRewardWallet(ctx context.Context, userID uuid.UUID, at time.Time) (usecase.RewardWallet, error) {
	response, err := c.rpc.GetRewardWallet(ctx, userAt(userID, at))
	return decode[usecase.RewardWallet](response, err)
}

func (c *Client) ProcessAction(ctx context.Context, command usecase.ProcessActionCommand) (usecase.ProcessActionResult, error) {
	request := &avitoshav1.ProcessActionRequest{
		UserId: command.UserID.String(), EventId: command.EventID.String(), ActionType: string(command.ActionType),
		EntityId: command.EntityID, Category: command.Category, MetadataJson: command.Metadata,
		OccurredAt: formatTime(command.OccurredAt), Now: formatTime(command.Now),
	}
	response, err := c.rpc.ProcessAction(ctx, request)
	return decode[usecase.ProcessActionResult](response, err)
}

func (c *Client) Subscribe(userID uuid.UUID) realtime.EventSubscription {
	ctx, cancel := context.WithCancel(context.Background())
	subscription := &interfaceSubscription{messages: make(chan []model.DomainEvent, 16), cancel: cancel}
	stream, err := c.rpc.SubscribeEvents(ctx, &avitoshav1.SubscribeEventsRequest{UserId: userID.String()})
	if err != nil {
		close(subscription.messages)
		return subscription
	}
	go subscription.receive(stream)
	return subscription
}

// interfaceSubscription is returned by value so every WebSocket owns one
// independent gRPC stream. Close remains idempotent through sync.Once.
type interfaceSubscription struct {
	messages chan []model.DomainEvent
	cancel   context.CancelFunc
	once     sync.Once
}

func (s *interfaceSubscription) Messages() <-chan []model.DomainEvent { return s.messages }
func (s *interfaceSubscription) Close()                               { s.once.Do(s.cancel) }

func (s *interfaceSubscription) receive(stream avitoshav1.GameService_SubscribeEventsClient) {
	defer close(s.messages)
	for {
		batch, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				s.cancel()
			}
			return
		}
		var events []model.DomainEvent
		if json.Unmarshal(batch.GetPayloadJson(), &events) != nil {
			s.cancel()
			return
		}
		select {
		case s.messages <- events:
		case <-stream.Context().Done():
			return
		}
	}
}

func decode[T any](response *avitoshav1.JsonResponse, err error) (T, error) {
	var value T
	if err != nil {
		return value, fmt.Errorf("game grpc request: %w", internalrpc.DecodeGameError(err))
	}
	if response == nil || json.Unmarshal(response.GetPayloadJson(), &value) != nil {
		return value, fmt.Errorf("decode game grpc response: %w", usecase.ErrUnexpectedStorage)
	}
	return value, nil
}

func userAt(userID uuid.UUID, at time.Time) *avitoshav1.UserAtRequest {
	return &avitoshav1.UserAtRequest{UserId: userID.String(), At: formatTime(at)}
}

func taskRequest(userID, taskID uuid.UUID, at time.Time) *avitoshav1.TaskRequest {
	return &avitoshav1.TaskRequest{UserId: userID.String(), TaskId: taskID.String(), At: formatTime(at)}
}

func formatTime(value time.Time) string { return value.UTC().Format(time.RFC3339Nano) }
