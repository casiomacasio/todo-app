package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/casiomacasio/todo-app/backend/internal/domain"
	"github.com/casiomacasio/todo-app/backend/internal/repository"
	"github.com/redis/go-redis/v9"
)

const (
	todoListCacheTTL = 10 * time.Minute
)

type TodoListService struct {
	repo        repository.TodoList
	redisClient *redis.Client
}

func NewTodoListService(repo repository.TodoList, redisClient *redis.Client) *TodoListService {
	return &TodoListService{repo: repo, redisClient: redisClient}
}

func (s *TodoListService) Create(ctx context.Context, userId int, list domain.CreateListRequest) (int, error) {
	id, err := s.repo.Create(userId, list)
	if err != nil {
		return 0, err
	}
	userKey := getUserListsCacheKey(userId)
	_ = s.redisClient.Del(ctx, userKey).Err() 
	return id, nil
}

func (s *TodoListService) GetAll(ctx context.Context, userId int) ([]domain.TodoList, error) {
	userKey := getUserListsCacheKey(userId)

	cached, err := s.redisClient.Get(ctx, userKey).Result()
	if err == nil {
		var lists []domain.TodoList
		if err := json.Unmarshal([]byte(cached), &lists); err == nil {
			return lists, nil
		}
	}

	lists, err := s.repo.GetAll(userId)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(lists)
	_ = s.redisClient.Set(ctx, userKey, data, todoListCacheTTL).Err()
	return lists, nil
}

func (s *TodoListService) GetById(ctx context.Context, userId, listId int) (domain.TodoList, error) {
	listKey := getListCacheKey(userId, listId)

	cached, err := s.redisClient.Get(ctx, listKey).Result()
	if err == nil {
		var list domain.TodoList
		if err := json.Unmarshal([]byte(cached), &list); err == nil {
			return list, nil
		}
	}

	list, err := s.repo.GetById(userId, listId)
	if err != nil {
		return domain.TodoList{}, err
	}

	data, _ := json.Marshal(list)
	_ = s.redisClient.Set(ctx, listKey, data, todoListCacheTTL).Err()
	return list, nil
}

func (s *TodoListService) UpdateById(ctx context.Context, userId, listId int, title, description *string) error {
	err := s.repo.UpdateById(userId, listId, title, description)
	if err != nil {
		return err
	}

	listKey := getListCacheKey(userId, listId)
	userKey := getUserListsCacheKey(userId)
	_ = s.redisClient.Del(ctx, listKey).Err()
	_ = s.redisClient.Del(ctx, userKey).Err()

	return nil
}

func (s *TodoListService) DeleteById(ctx context.Context, userId, listId int) error {
	err := s.repo.DeleteById(userId, listId)
	if err != nil {
		return err
	}

	listKey := getListCacheKey(userId, listId)
	userKey := getUserListsCacheKey(userId)
	_ = s.redisClient.Del(ctx, listKey).Err()
	_ = s.redisClient.Del(ctx, userKey).Err()

	return nil
}

func getUserListsCacheKey(userId int) string {
	return fmt.Sprintf("user:%d:lists", userId)
}

func getListCacheKey(userID, listID int) string {
	return fmt.Sprintf("user:%d:list:%d", userID, listID)
}