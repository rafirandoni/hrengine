package hrengine

import (
	"github.com/rafirandoni/hrengine/server"
	"github.com/rafirandoni/hrengine/utility"
	"github.com/rafirandoni/hrengine/web"
)

var (
	Env    utility.Env
	Server server.Server
	Web    web.Web
)

func init() {
	Env = utility.NewEnv()
	Server = server.New()
	Web = web.New()
}
