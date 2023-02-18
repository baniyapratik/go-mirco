package sever

type Config struct {
	// Host on which the service will be enabled
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	// Port on which the service will be enabled
	WebPort int `env:"WEB_PORT" envDefault:"8083"`
	// Port on which the service will be enabled
	RPCPort                   int    `env:"RPC_PORT" envDefault:"5001"`
	GRPCPort                  int    `env:"GRPC_PORT" envDefault:"50001"`
	MongoURL                  string `env:"MONGO_URL" envDefault:"mongodb://mongo:27017"`
	ServiceDatabaseName       string `env:"SERVICE_DATABASE_NAME" envDefault:"logger-service"`
	ServiceLogsCollectionName string `env:"SERVICE_LOGS_COLLECTION_NAME" envDefault:"service-logs"`
}
