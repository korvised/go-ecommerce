package mytests

import (
	"encoding/json"
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/servers"
	"github.com/korvised/go-ecommerce/pkg/databases"
)

func SetupTest() servers.IModuleFactory {
	cfg := config.LoadConfig("../.env.test")

	db := databases.DbConnect(cfg.Db())

	s := servers.NewServer(cfg, db)
	return servers.InitModule(nil, s.GetServer(), nil)
}

func CompressToJson(obj any) string {
	result, _ := json.Marshal(&obj)
	return string(result)
}
