package mongo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	MongoClient   *mongo.Client
	MongoDatabase *mongo.Database
	logger        *slog.Logger
}

type MongoConn struct {
	host               string
	port               string
	user               string
	password           string
	dbName             string
	authDb             string
	withCredentials    bool
	connectTimeout     time.Duration
	connectionAttempts int
}

type connopt func(*MongoConn)

const (
	defaultConnAttempts   = 10
	defaultConnectTimeout = 2 * time.Second
)

func NewMongoConn(
	host, port, dbName string,
	connopts ...connopt) *MongoConn {
	mongoConn := &MongoConn{
		host:               host,
		port:               port,
		user:               "",
		password:           "",
		withCredentials:    false,
		connectTimeout:     defaultConnectTimeout,
		dbName:             dbName,
		authDb:             "",
		connectionAttempts: defaultConnAttempts,
	}

	for _, opt := range connopts {
		opt(mongoConn)
	}

	return mongoConn

}

func WithTimeout(connTimeout time.Duration) connopt {
	return func(mongoConn *MongoConn) {
		mongoConn.connectTimeout = connTimeout
	}
}

func WithCredentials(user, password string) connopt {
	return func(mongoConn *MongoConn) {
		mongoConn.user = user
		mongoConn.password = password
		mongoConn.authDb = mongoConn.dbName
		mongoConn.withCredentials = true
	}
}

func WithAuthDb(authDb string) connopt {
	return func(mongoConn *MongoConn) {
		mongoConn.authDb = authDb
	}
}

func WithConnectionAttempts(connectAttempts int) connopt {
	return func(mongoConn *MongoConn) {
		mongoConn.connectionAttempts = connectAttempts
	}
}

func (mongoConn *MongoConn) parseUrl() string {
	if mongoConn.withCredentials {
		return fmt.Sprintf(
			"mongodb://%s:%s@%s:%s/%s",
			mongoConn.user,
			mongoConn.password,
			mongoConn.host,
			mongoConn.port,
			mongoConn.dbName,
		)
	} else {
		return fmt.Sprintf(
			"mongodb://%s:%s",
			mongoConn.host,
			mongoConn.port,
		)
	}
}

func New(ctx context.Context, mongoConn *MongoConn, logger *slog.Logger) (*MongoClient, error) {
	mongoClient := &MongoClient{
		logger: logger,
	}

	cont, cancel := context.WithTimeout(context.Background(), mongoConn.connectTimeout)
	defer cancel()

	url := mongoConn.parseUrl()
	opts := options.Client().ApplyURI(url)
	fmt.Println(url)
	if mongoConn.withCredentials {
		opts.SetAuth(options.Credential{
			AuthSource: mongoConn.authDb,
			Username:   mongoConn.user,
			Password:   mongoConn.password,
		})
	}

	client, err := mongo.Connect(cont, opts)
	if err != nil {
		return nil, fmt.Errorf("Mongo - New - Connect: %w", err)
	}
	mongoClient.MongoClient = client
	if err := mongoClient.pingWithAttempts(cont, mongoConn.connectionAttempts); err != nil {
		return nil, err
	}

	mongoClient.MongoDatabase = client.Database(mongoConn.dbName)
	return mongoClient, nil
}

func (mongoClient *MongoClient) pingWithAttempts(ctx context.Context, attempts int) error {
	var err error
	for attempts > 0 {
		attempts--
		if err = mongoClient.MongoClient.Ping(ctx, nil); err == nil {
			return nil
		}
		mongoClient.logger.Warn("Ping to mongo is failed, ", slog.Int("attempts", attempts))
	}

	return fmt.Errorf("Mongo - PingWithAttempts: Zero attempts left, failed to ping mongo: %w", err)
}
