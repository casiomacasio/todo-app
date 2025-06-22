package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/casiomacasio/todo-app/internal/repository"
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

func (s *TodoListService) Create(userId int, list domain.TodoList) (int, error) {
	id, err := s.repo.Create(userId, list)
	if err != nil {
		return 0, err
	}
	userKey := getUserListsCacheKey(userId)
	s.redisClient.Del(ctx, userKey)
	return id, nil
}

func (s *TodoListService) GetAll(userId int) ([]domain.TodoList, error) {
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
	s.redisClient.Set(ctx, userKey, data, todoListCacheTTL)
	return lists, nil
}

func (s *TodoListService) GetById(userId, listId int) (domain.TodoList, error) {
	listKey := getListCacheKey(userId)

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
	s.redisClient.Set(ctx, listKey, data, todoListCacheTTL)
	return list, nil
}

func (s *TodoListService) UpdateById(userId, listId int, title, description *string) error {
	err := s.repo.UpdateById(userId, listId, title, description)
	if err != nil {
		return err
	}

	listKey := getListCacheKey(listId)
	userKey := getUserListsCacheKey(userId)
	s.redisClient.Del(ctx, listKey)
	s.redisClient.Del(ctx, userKey)

	return nil
}

func (s *TodoListService) DeleteById(userId, listId int) error {
	err := s.repo.DeleteById(userId, listId)
	if err != nil {
		return err
	}

	listKey := getListCacheKey(listId)
	userKey := getUserListsCacheKey(userId)
	s.redisClient.Del(ctx, listKey)
	s.redisClient.Del(ctx, userKey)

	return nil
}

func getUserListsCacheKey(userId int) string {
	return fmt.Sprintf("user:%d:lists", userId)
}

func getListCacheKey(listId int) string {
	return fmt.Sprintf("list:%d", listId)
}
