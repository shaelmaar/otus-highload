package postfeed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/cache"
	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

const (
	defaultTTL                  = 10 * time.Minute
	defaultLockTimeout          = time.Second
	defaultWaitForCacheDuration = 500 * time.Millisecond
	defaultCheckCacheInterval   = 100 * time.Millisecond
	defaultReadTimeout          = 500 * time.Millisecond
	lockSuffix                  = ":lock"
)

type Cache struct {
	provider cache.ValkeyProvider
	logger   *zap.Logger

	ttl                  time.Duration
	lockTimeout          time.Duration
	waitForCacheDuration time.Duration
	cacheCheckInterval   time.Duration
	readTimeout          time.Duration
}

func New(
	provider cache.ValkeyProvider,
	logger *zap.Logger,
) (*Cache, error) {
	if utils.IsNil(provider) {
		return nil, errors.New("provider is nil")
	}

	if utils.IsNil(logger) {
		return nil, errors.New("logger is nil")
	}

	if utils.IsNil(provider) {
		return nil, errors.New("provider is nil")
	}

	return &Cache{
		provider: provider,
		logger:   logger,

		ttl:                  defaultTTL,
		lockTimeout:          defaultLockTimeout,
		waitForCacheDuration: defaultWaitForCacheDuration,
		cacheCheckInterval:   defaultCheckCacheInterval,
		readTimeout:          defaultReadTimeout,
	}, nil
}

func (c *Cache) GetUserFeed(ctx context.Context, input dto.GetPostFeed) ([]domain.Post, bool, error) {
	client, err := c.provider.Client()
	if err != nil {
		c.logger.Error("failed to get client", zap.Error(err))

		return nil, false, nil
	}

	key := postFeedCacheKey(input.UserID)

	start := int64(input.Offset)
	stop := start + int64(input.Limit) - 1

	readCtx, readCancel := context.WithTimeout(ctx, c.readTimeout)
	defer readCancel()

	posts, got, err := c.getUserFeed(readCtx, client, key, start, stop)
	if err != nil {
		c.provider.ResetClient()

		c.logger.Error("failed to get posts from valkey", zap.Error(err))

		return nil, false, fmt.Errorf("failed to get posts from valkey: %w", err)
	}

	if got {
		return posts, true, nil
	}

	if !c.lockExists(readCtx, client, key) {
		return nil, false, nil
	}

	waitCtx, waitCancel := context.WithTimeout(ctx, c.waitForCacheDuration)
	defer waitCancel()

	ticker := time.NewTicker(c.waitForCacheDuration)

	for {
		select {
		case <-ticker.C:
			posts, got, err = c.getUserFeed(waitCtx, client, key, start, stop)
			if err != nil {
				return nil, false, fmt.Errorf("failed to get posts from valkey: %w", err)
			}

			if got {
				return posts, true, nil
			}
		case <-waitCtx.Done():
			return nil, false, nil
		}
	}
}

func (c *Cache) SetUserFeed(ctx context.Context, userID uuid.UUID, posts []domain.Post) error {
	client, err := c.provider.Client()
	if err != nil {
		c.logger.Error("failed to get client", zap.Error(err))

		return nil
	}

	key := postFeedCacheKey(userID)

	locked, err := c.lock(ctx, client, key)
	if err != nil {
		return fmt.Errorf("failed to lock post feed: %w", err)
	}

	if !locked {
		return nil
	}

	//nolint:contextcheck // тут нужен отдельный контекст.
	defer func() {
		err = c.unlock(context.Background(), client, key)
		if err != nil {
			c.logger.Error("failed to unlock post feed", zap.Error(err))
		}
	}()

	var cmds []valkey.Completed

	cmds = append(cmds, client.B().Del().Key(key).Build())

	jsonData, err := json.Marshal(posts)
	if err != nil {
		return fmt.Errorf("failed to serialize posts: %w", err)
	}

	cmds = append(cmds, client.B().Set().Key(key).Value(string(jsonData)).Ex(c.ttl).Build())

	for _, resp := range client.DoMulti(ctx, cmds...) {
		if err := resp.Error(); err != nil {
			return fmt.Errorf("failed to set post feed: %w", err)
		}
	}

	return nil
}

func (c *Cache) UserFeedExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	client, err := c.provider.Client()
	if err != nil {
		c.logger.Error("failed to get client", zap.Error(err))

		return false, nil
	}

	key := postFeedCacheKey(userID)

	_, err = client.DoCache(ctx, client.B().Get().Key(key).Cache(), c.ttl).ToString()

	switch {
	case valkey.IsValkeyNil(err):
		return false, nil
	case err != nil:
		c.provider.ResetClient()

		return false, fmt.Errorf("failed to get posts from valkey: %w", err)
	}

	return true, nil
}

func (c *Cache) getUserFeed(
	ctx context.Context, client valkey.Client, key string, start, stop int64) ([]domain.Post, bool, error) {
	var posts []domain.Post

	cmd := client.B().Get().Key(key).Cache()
	jsonData, err := client.DoCache(ctx, cmd, c.ttl).ToString()

	switch {
	case valkey.IsValkeyNil(err):
		return nil, false, nil
	case err != nil:
		c.provider.ResetClient()

		return nil, false, fmt.Errorf("failed to get posts from valkey: %w", err)
	}

	err = json.Unmarshal([]byte(jsonData), &posts)
	if err != nil {
		return nil, false, fmt.Errorf("failed to deserialize posts: %w", err)
	}

	return utils.SafeSliceRange(posts, int(start), int(stop)), true, nil
}

func (c *Cache) lockExists(ctx context.Context, client valkey.Client, key string) bool {
	exists, _ := client.Do(
		ctx,
		client.B().Exists().Key(lockFeedCacheKey(key)).Build(),
	).ToInt64()

	return exists > 0
}

func (c *Cache) lock(ctx context.Context, client valkey.Client, key string) (bool, error) {
	cmd := client.B().Set().Key(lockFeedCacheKey(key)).Value("locked").Nx().Px(c.lockTimeout).Build()

	err := client.Do(ctx, cmd).Error()
	if err != nil {
		if errors.Is(err, valkey.Nil) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (c *Cache) unlock(ctx context.Context, client valkey.Client, key string) error {
	cmd := client.B().Del().Key(lockFeedCacheKey(key)).Build()

	err := client.Do(ctx, cmd).Error()
	if err != nil {
		if errors.Is(err, valkey.Nil) {
			return nil
		}

		return err
	}

	return nil
}

func postFeedCacheKey(userID uuid.UUID) string {
	return fmt.Sprintf("post_feed:%s", userID.String())
}

func lockFeedCacheKey(key string) string {
	return fmt.Sprintf("%s%s", key, lockSuffix)
}
