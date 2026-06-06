package user

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewRepository,
	wire.Bind(new(Repository), new(*repository)),
	NewService,
	wire.Bind(new(Service), new(*service)),
)
