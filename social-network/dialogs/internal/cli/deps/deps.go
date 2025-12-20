package deps

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/config"
	httpServer "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/httptransport/server"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/kafka/consumer/countersmessages"
)

func (c *Container) Config() *config.Config {
	return do.MustInvoke[*config.Config](c.i)
}

func (c *Container) Logger() *zap.Logger {
	return do.MustInvoke[*zap.Logger](c.i)
}

func (c *Container) HTTPServer() *httpServer.Server {
	return do.MustInvokeNamed[*httpServer.Server](c.i, nameHTTPServer)
}

func (c *Container) DebugServer() *echo.Echo {
	return do.MustInvokeNamed[*echo.Echo](c.i, nameDebugServer)
}

func (c *Container) CountersMessagesConsumer() *countersmessages.Consumer {
	return do.MustInvokeNamed[*countersmessages.Consumer](c.i, nameCountersMessagesConsumer)
}
