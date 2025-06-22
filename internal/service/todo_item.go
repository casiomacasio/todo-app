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
	todoCacheTTL = 10 * time.Minute
)

type TodoItemService struct {
	repo        repository.TodoItem
	redisClient *redis.Client
}

func NewTodoItemService(repo repository.TodoItem, redisClient *redis.Client) *TodoItemService {
	return &TodoItemService{repo: repo, redisClient: redisClient}
}

func (s *TodoItemService) Create(userId, listId int, item domain.TodoItem) (int, error) {
	id, err := s.repo.Create(userId, listId, item)
	if err != nil {
		return 0, err
	}
	listKey := getListItemsCacheKey(userId, listId)
	s.redisClient.Del(ctx, listKey)
	return id, nil
}

func (s *TodoItemService) GetAllItems(userId, listId int) ([]domain.TodoItem, error) {
	listKey := getListItemsCacheKey(userId, listId)

	cached, err := s.redisClient.Get(ctx, listKey).Result()
	if err == nil {
		var items []domain.TodoItem
		if err := json.Unmarshal([]byte(cached), &items); err == nil {
			return items, nil
		}
	}

	items, err := s.repo.GetAllItems(userId, listId)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(items)
	s.redisClient.Set(ctx, listKey, data, todoCacheTTL)
	return items, nil
}

func (s *TodoItemService) GetById(userId, itemId int) (domain.TodoItem, error) {
	itemKey := getItemCacheKey(itemId)

	cached, err := s.redisClient.Get(ctx, itemKey).Result()
	if err == nil {
		var item domain.TodoItem
		if err := json.Unmarshal([]byte(cached), &item); err == nil {
			return item, nil
		}
	}

	item, err := s.repo.GetById(userId, itemId)
	if err != nil {
		return domain.TodoItem{}, err
	}

	data, _ := json.Marshal(item)
	s.redisClient.Set(ctx, itemKey, data, todoCacheTTL)
	return item, nil
}

func (s *TodoItemService) UpdateById(userId, itemId int, title, description *string, done *bool) error {
	err := s.repo.UpdateById(userId, itemId, title, description, done)
	if err != nil {
		return err
	}

	itemKey := getItemCacheKey(itemId)
	s.redisClient.Del(ctx, itemKey)

	return nil
}

func (s *TodoItemService) DeleteById(userId, itemId int) error {
	err := s.repo.DeleteById(userId, itemId)
	if err != nil {
		return err
	}

	itemKey := getItemCacheKey(itemId)
	s.redisClient.Del(ctx, itemKey)

	return nil
}

func getListItemsCacheKey(userId, listId int) string {
	return fmt.Sprintf("user:%d:list:%d:items", userId, listId)
}

func getItemCacheKey(itemId int) string {
	return fmt.Sprintf("todo_item:%d", itemId)
}
