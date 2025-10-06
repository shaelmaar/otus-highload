package deps

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/config"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/server"
	"github.com/shaelmaar/otus-highload/social-network/internal/rabbitmq"
	userUseCases "github.com/shaelmaar/otus-highload/social-network/internal/usecase/user"
)

func (c *Container) Config() *config.Config {
	return do.MustInvoke[*config.Config](c.i)
}

func (c *Container) PgxPool() *pgxpool.Pool {
	return do.MustInvokeNamed[*pgxpool.Pool](c.i, namePgxPool)
}

func (c *Container) Logger() *zap.Logger {
	return do.MustInvoke[*zap.Logger](c.i)
}

func (c *Container) HTTPServer() *server.Server {
	return do.MustInvoke[*server.Server](c.i)
}

func (c *Container) DebugServer() *echo.Echo {
	return do.MustInvokeNamed[*echo.Echo](c.i, nameDebugServer)
}

func (c *Container) UserUseCases() *userUseCases.UseCases {
	return do.MustInvoke[*userUseCases.UseCases](c.i)
}

func (c *Container) UserFeedTaskConsumer() *rabbitmq.Consumer[dto.UserFeedUpdateTask] {
	return do.MustInvokeNamed[*rabbitmq.Consumer[dto.UserFeedUpdateTask]](c.i, nameUserFeedTaskConsumer)
}

func (c *Container) UserFeedChunkedTaskConsumer() *rabbitmq.Consumer[dto.UserFeedChunkedUpdateTask] {
	return do.MustInvokeNamed[*rabbitmq.Consumer[dto.UserFeedChunkedUpdateTask]](c.i, nameUserFeedChunkedTaskConsumer)
}
