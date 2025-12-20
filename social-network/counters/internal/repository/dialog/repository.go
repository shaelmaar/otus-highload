package dialog

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/valkey-io/valkey-go"

	"github.com/shaelmaar/otus-highload/social-network/counters/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/counters/pkg/utils"
)

type Repository struct {
	client          valkey.Client
	_replicaClients []valkey.Client
}

func New(client valkey.Client, replicaClients []valkey.Client) (*Repository, error) {
	if utils.IsNil(client) {
		return nil, errors.New("client is nil")
	}

	if utils.IsNil(replicaClients) {
		return nil, errors.New("replica clients is nil")
	}

	return &Repository{
		client:          client,
		_replicaClients: replicaClients,
	}, nil
}

func (r *Repository) CountUnreadMessages(ctx context.Context, key domain.UnreadDialogMessageCountKey) (int64, error) {
	cmd := r.client.B().Get().Key(key.String()).Build()

	val, err := r.client.Do(ctx, cmd).ToString()

	switch {
	case errors.Is(err, valkey.Nil):
		return 0, nil
	case err != nil:
		return 0, fmt.Errorf("failed to load unread dialog message count from valkey: %w", err)
	}

	res, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse unread dialog message count from valkey: %w", err)
	}

	return res, nil
}

func (r *Repository) IncrementUnreadMessages(
	ctx context.Context,
	key domain.UnreadDialogMessageCountKey,
	idempotencyKey string,
) error {
	exists, _ := r.checkIdempotencyKey(ctx, idempotencyKey)
	if exists {
		return nil
	}

	cmd := r.client.B().Incr().Key(key.String()).Build()

	err := r.client.Do(ctx, cmd).Error()

	switch {
	case errors.Is(err, valkey.Nil):
		return nil
	case err != nil:
		return fmt.Errorf("failed to increment unread dialog messages in valkey: %w", err)
	}

	_ = r.setIdempotencyKey(ctx, idempotencyKey)

	return nil
}

func (r *Repository) DecrementUnreadMessages(
	ctx context.Context,
	key domain.UnreadDialogMessageCountKey,
	decrBy int,
	idempotencyKey string,
) error {
	exists, _ := r.checkIdempotencyKey(ctx, idempotencyKey)
	if exists {
		return nil
	}

	defer func() {
		_ = r.setIdempotencyKey(context.Background(), idempotencyKey)
	}()

	locked, err := r.lock(ctx, key.String())
	if err != nil {
		return fmt.Errorf("failed to lock unread dialog messages in valkey: %w", err)
	}

	if !locked {
		return nil
	}

	defer func() {
		_ = r.unlock(context.Background(), key.String())
	}()

	cmd := r.client.B().Get().Key(key.String()).Build()

	val, err := r.client.Do(ctx, cmd).ToString()

	switch {
	case errors.Is(err, valkey.Nil):
		return nil
	case err != nil:
		return fmt.Errorf("failed to get unread dialog messages in valkey: %w", err)
	}

	counter, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse unread dialog messages counter in valkey: %w", err)
	}

	cmd = r.client.B().Decrby().Key(key.String()).Decrement(int64(decrBy)).Build()

	if int(counter)-decrBy <= 0 {
		cmd = r.client.B().Set().Key(key.String()).Value("0").Build()
	}

	err = r.client.Do(ctx, cmd).Error()
	if err != nil {
		return fmt.Errorf("failed to decrement unread dialog messages in valkey: %w", err)
	}

	return nil
}

func (r *Repository) Slave() domain.DialogSlaveRepository {
	if len(r._replicaClients) == 0 {
		return &Repository{
			client:          r.client,
			_replicaClients: nil,
		}
	}

	replicaClient := r._replicaClients[rand.Intn(len(r._replicaClients))]

	return &Repository{
		client:          replicaClient,
		_replicaClients: nil,
	}
}

func (r *Repository) checkIdempotencyKey(ctx context.Context, key string) (bool, error) {
	if key == "" {
		return false, nil
	}

	cmd := r.client.B().Exists().Key(key).Build()

	exists, err := r.client.Do(ctx, cmd).ToBool()
	if err != nil {
		return false, fmt.Errorf("failed to check if key exists in valkey: %w", err)
	}

	return exists, nil
}

func (r *Repository) setIdempotencyKey(ctx context.Context, key string) error {
	if key == "" {
		return nil
	}

	cmd := r.client.B().Set().Key(key).Value("1").Ex(5 * time.Minute).Build()

	err := r.client.Do(ctx, cmd).Error()
	if err != nil {
		return fmt.Errorf("failed to set key in valkey: %w", err)
	}

	return nil
}

func (r *Repository) lock(ctx context.Context, key string) (bool, error) {
	cmd := r.client.B().Set().Key(lockKey(key)).Value("locked").Nx().Px(time.Second).Build()

	err := r.client.Do(ctx, cmd).Error()
	if err != nil {
		if errors.Is(err, valkey.Nil) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *Repository) unlock(ctx context.Context, key string) error {
	cmd := r.client.B().Del().Key(lockKey(key)).Build()

	err := r.client.Do(ctx, cmd).Error()
	if err != nil {
		if errors.Is(err, valkey.Nil) {
			return nil
		}

		return err
	}

	return nil
}

func lockKey(key string) string {
	return fmt.Sprintf("lock:%s", key)
}
