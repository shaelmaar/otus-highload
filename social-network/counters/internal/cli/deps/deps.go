package deps

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"github.com/valkey-io/valkey-go"
	"go.uber.org/zap"

	httpServer "github.com/shaelmaar/otus-highload/social-network/counters/internal/httptransport/server"
	"github.com/shaelmaar/otus-highload/social-network/counters/internal/kafka/consumer/dialogsmessages"

	"github.com/shaelmaar/otus-highload/social-network/counters/internal/config"
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

func (c *Container) DialogsMessagesConsumer() *dialogsmessages.Consumer {
	return do.MustInvokeNamed[*dialogsmessages.Consumer](c.i, nameDialogsMessagesConsumer)
}

func (c *Container) DebugServer() *echo.Echo {
	return do.MustInvokeNamed[*echo.Echo](c.i, nameDebugServer)
}

// todo потом удалить.
func (c *Container) ValkeyMasterClient() valkey.Client {
	return do.MustInvokeNamed[valkey.Client](c.i, nameValkeyMasterClient)
}
