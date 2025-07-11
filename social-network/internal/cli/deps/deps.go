package deps

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/config"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/server"
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
