package gameserver

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	avitoshav1 "github.com/guitaramust-sudo/Avitosha/app/backend/internal/gen/avitosha/v1"
	"github.com/guitaramust-sudo/Avitosha/app/backend/internal/model"
	"github.com/guitaramust-sudo/Avitosha/app/backend/internal/realtime"
	internalrpc "github.com/guitaramust-sudo/Avitosha/app/backend/internal/rpc"
	"github.com/guitaramust-sudo/Avitosha/app/backend/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GameUseCase interface {
	EnsureProfile(context.Context, uuid.UUID, time.Time) (usecase.GameProfile, error)
	RenamePet(context.Context, uuid.UUID, string, time.Time) (usecase.GameProfile, error)
	ListTasks(context.Context, uuid.UUID, time.Time) ([]model.TaskProgress, error)
	GetTask(context.Context, uuid.UUID, uuid.UUID, time.Time) (model.TaskProgress, error)
	GetTaskAdvice(context.Context, uuid.UUID, uuid.UUID, time.Time) (usecase.TaskAdvice, error)
	GetRoom(context.Context, uuid.UUID, time.Time) ([]model.RoomItemProgress, error)
	GetStory(context.Context, uuid.UUID, time.Time) (model.StorySnapshot, error)
	GetDailySummary(context.Context, uuid.UUID, time.Time) (usecase.DailySummary, error)
	GetLeaderboard(context.Context, uuid.UUID, int, time.Time) (usecase.Leaderboard, error)
	GetAchievements(context.Context, uuid.UUID, time.Time) ([]model.AchievementProgress, error)
	GetRewardBalances(context.Context, uuid.UUID, time.Time) ([]model.RewardBalance, error)
	GetRewardWallet(context.Context, uuid.UUID, time.Time) (usecase.RewardWallet, error)
	ProcessAction(context.Context, usecase.ProcessActionCommand) (usecase.ProcessActionResult, error)
}

type Server struct {
	avitoshav1.UnimplementedGameServiceServer
	game GameUseCase
	hub  *realtime.Hub
}

func New(game GameUseCase, hub *realtime.Hub) *Server { return &Server{game: game, hub: hub} }

func (s *Server) EnsureProfile(ctx context.Context, request *avitoshav1.UserAtRequest) (*avitoshav1.JsonResponse, error) {
	userID, at, err := parseUserAt(request.GetUserId(), request.GetAt())
	if err != nil {
		return nil, err
	}
	return encode(s.game.EnsureProfile(ctx, userID, at))
}

func (s *Server) RenamePet(ctx context.Context, request *avitoshav1.RenamePetRequest) (*avitoshav1.JsonResponse, error) {
	userID, at, err := parseUserAt(request.GetUserId(), request.GetAt())
	if err != nil {
		return nil, err
	}
	return encode(s.game.RenamePet(ctx, userID, request.GetName(), at))
}

func (s *Server) ListTasks(ctx context.Context, request *avitoshav1.UserAtRequest) (*avitoshav1.JsonResponse, error) {
	userID, at, err := parseUserAt(request.GetUserId(), request.GetAt())
	if err != nil {
		return nil, err
	}
	return encode(s.game.ListTasks(ctx, userID, at))
}

func (s *Server) GetTask(ctx context.Context, request *avitoshav1.TaskRequest) (*avitoshav1.JsonResponse, error) {
	userID, taskID, at, err := parseTaskRequest(request)
	if err != nil {
		return nil, err
	}
	return encode(s.game.GetTask(ctx, userID, taskID, at))
}

func (s *Server) GetTaskAdvice(ctx context.Context, request *avitoshav1.TaskRequest) (*avitoshav1.JsonResponse, error) {
	userID, taskID, at, err := parseTaskRequest(request)
	if err != nil {
		return nil, err
	}
	return encode(s.game.GetTaskAdvice(ctx, userID, taskID, at))
}

func (s *Server) GetRoom(ctx context.Context, request *avitoshav1.UserAtRequest) (*avitoshav1.JsonResponse, error) {
	userID, at, err := parseUserAt(request.GetUserId(), request.GetAt())
	if err != nil {
		return nil, err
	}
	return encode(s.game.GetRoom(ctx, userID, at))
}

func (s *Server) GetStory(ctx context.Context, request *avitoshav1.UserAtRequest) (*avitoshav1.JsonResponse, error) {
	userID, at, err := parseUserAt(request.GetUserId(), request.GetAt())
	if err != nil {
		return nil, err
	}
	return encode(s.game.GetStory(ctx, userID, at))
}

func (s *Server) GetDailySummary(ctx context.Context, request *avitoshav1.UserAtRequest) (*avitoshav1.JsonResponse, error) {
	userID, at, err := parseUserAt(request.GetUserId(), request.GetAt())
	if err != nil {
		return nil, err
	}
	return encode(s.game.GetDailySummary(ctx, userID, at))
}

func (s *Server) GetLeaderboard(ctx context.Context, request *avitoshav1.LeaderboardRequest) (*avitoshav1.JsonResponse, error) {
	userID, at, err := parseUserAt(request.GetUserId(), request.GetAt())
	if err != nil {
		return nil, err
	}
	return encode(s.game.GetLeaderboard(ctx, userID, int(request.GetLimit()), at))
}

func (s *Server) GetAchievements(ctx context.Context, request *avitoshav1.UserAtRequest) (*avitoshav1.JsonResponse, error) {
	userID, at, err := parseUserAt(request.GetUserId(), request.GetAt())
	if err != nil {
		return nil, err
	}
	return encode(s.game.GetAchievements(ctx, userID, at))
}

func (s *Server) GetRewardBalances(ctx context.Context, request *avitoshav1.UserAtRequest) (*avitoshav1.JsonResponse, error) {
	userID, at, err := parseUserAt(request.GetUserId(), request.GetAt())
	if err != nil {
		return nil, err
	}
	return encode(s.game.GetRewardBalances(ctx, userID, at))
}

func (s *Server) GetRewardWallet(ctx context.Context, request *avitoshav1.UserAtRequest) (*avitoshav1.JsonResponse, error) {
	userID, at, err := parseUserAt(request.GetUserId(), request.GetAt())
	if err != nil {
		return nil, err
	}
	return encode(s.game.GetRewardWallet(ctx, userID, at))
}

func (s *Server) ProcessAction(ctx context.Context, request *avitoshav1.ProcessActionRequest) (*avitoshav1.JsonResponse, error) {
	userID, err := parseUUID("user_id", request.GetUserId())
	if err != nil {
		return nil, err
	}
	eventID, err := parseUUID("event_id", request.GetEventId())
	if err != nil {
		return nil, err
	}
	occurredAt, err := parseTime("occurred_at", request.GetOccurredAt())
	if err != nil {
		return nil, err
	}
	now, err := parseTime("now", request.GetNow())
	if err != nil {
		return nil, err
	}
	command := usecase.ProcessActionCommand{
		UserID: userID, EventID: eventID, ActionType: model.ActionType(request.GetActionType()),
		EntityID: request.EntityId, Category: request.Category, Metadata: json.RawMessage(request.GetMetadataJson()),
		OccurredAt: occurredAt, Now: now,
	}
	return encode(s.game.ProcessAction(ctx, command))
}

func (s *Server) SubscribeEvents(request *avitoshav1.SubscribeEventsRequest, stream avitoshav1.GameService_SubscribeEventsServer) error {
	userID, err := parseUUID("user_id", request.GetUserId())
	if err != nil {
		return err
	}
	subscription := s.hub.Subscribe(userID)
	defer subscription.Close()
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case events, ok := <-subscription.Messages():
			if !ok {
				return nil
			}
			payload, marshalErr := json.Marshal(events)
			if marshalErr != nil {
				return status.Error(codes.Internal, "encode event batch")
			}
			if sendErr := stream.Send(&avitoshav1.EventBatch{PayloadJson: payload}); sendErr != nil {
				return sendErr
			}
		}
	}
}

func encode[T any](value T, err error) (*avitoshav1.JsonResponse, error) {
	if err != nil {
		return nil, internalrpc.GameError(err)
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, status.Error(codes.Internal, "encode response")
	}
	return &avitoshav1.JsonResponse{PayloadJson: payload}, nil
}

func parseUserAt(userIDValue, atValue string) (uuid.UUID, time.Time, error) {
	userID, err := parseUUID("user_id", userIDValue)
	if err != nil {
		return uuid.Nil, time.Time{}, err
	}
	at, err := parseTime("at", atValue)
	return userID, at, err
}

func parseTaskRequest(request *avitoshav1.TaskRequest) (uuid.UUID, uuid.UUID, time.Time, error) {
	userID, at, err := parseUserAt(request.GetUserId(), request.GetAt())
	if err != nil {
		return uuid.Nil, uuid.Nil, time.Time{}, err
	}
	taskID, err := parseUUID("task_id", request.GetTaskId())
	return userID, taskID, at, err
}

func parseUUID(field, value string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, status.Errorf(codes.InvalidArgument, "%s must be a UUID", field)
	}
	return parsed, nil
}

func parseTime(field, value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, status.Errorf(codes.InvalidArgument, "%s must be RFC3339", field)
	}
	return parsed, nil
}
