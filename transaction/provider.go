package transaction

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewRepository,
	wire.Bind(new(Repository), new(*repository)),
	NewTxRunner,
	wire.Bind(new(TxRunner), new(*gormTxRunner)),
	NewService,
	wire.Bind(new(Service), new(*service)),
)
