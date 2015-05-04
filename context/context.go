package context

import (
    "github.com/cosmos-io/cosmos/service"
)

type CosmosContext struct {
    CosmosService *service.CosmosService
    Params map[string]string
}
