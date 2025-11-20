package deps

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/config"
	grpcServer "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/grpctransport/server"
	httpServer "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/httptransport/server"
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

func (c *Container) GRPCServer() *grpcServer.Server {
	return do.MustInvokeNamed[*grpcServer.Server](c.i, nameGRPCServer)
}

func (c *Container) DebugServer() *echo.Echo {
	return do.MustInvokeNamed[*echo.Echo](c.i, nameDebugServer)
}
