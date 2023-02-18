package mongo_driver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"strings"
)

type Interactor interface {
	Find(ctx context.Context, d bson.D, opts *options.FindOptions) (*mongo.Cursor, error)
	Disconnect(ctx context.Context) error
	Insert(ctx context.Context, data interface{}) error
	Replace(ctx context.Context, key string, data interface{}) error
}

type Service struct {
	dbURL          string
	username       string
	password       string
	databaseName   string
	collectionName string
	caCertFile     string
	logger         log.Logger
	collection     *mongo.Collection
	client         *mongo.Client
}

func New(opts ...Option) (*Service, error) {
	s := Service{}
	for _, o := range opts {
		o(&s)
	}
	if err := s.validate(); err != nil {
		return nil, fmt.Errorf("error in validation %w", err)
	}
	if err := s.connect(context.TODO()); err != nil {
		return nil, fmt.Errorf("error in setting up connection %w", err)
	}
	return &s, nil
}
func (s *Service) Disconnect(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
func (s *Service) validate() error {
	var errMsg []string
	if s.dbURL == "" {
		errMsg = append(errMsg, "missing connection string")
	}
	if s.username == "" {
		errMsg = append(errMsg, "missing username for db connection")
	}
	if s.password == "" {
		errMsg = append(errMsg, "missing password for db connection")
	}
	if s.databaseName == "" {
		errMsg = append(errMsg, "missing database name for db connection")
	}
	if s.collectionName == "" {
		errMsg = append(errMsg, "missing collection name for db connection")
	}
	if s.caCertFile == "" {
		errMsg = append(errMsg, "missing ca cert for db connection")
	}
	if len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, ","))
	}
	return nil
}

func (s *Service) connect(ctx context.Context) error {
	credential := options.Credential{
		Username: s.username,
		Password: s.password,
	}

	caCert, err := ioutil.ReadFile(s.caCertFile)
	if err != nil {
		return fmt.Errorf("failed to read ca cert file %w", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	clientOpts := options.Client().ApplyURI(s.dbURL).SetAuth(credential).SetTLSConfig(tlsConfig).SetRetryWrites(false)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to connect to mongo DB %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping mongo DB %w", err)
	}
	s.collection = client.Database(s.databaseName).Collection(s.collectionName)
	s.client = client
	return nil
}

func (s *Service) Insert(ctx context.Context, data interface{}) error {
	_, err := s.collection.InsertOne(ctx, data)
	if err != nil {

	}
	return nil
}

func (s *Service) Find(ctx context.Context, d bson.D, opts *options.FindOptions) (*mongo.Cursor, error) {
	cursor, err := s.collection.Find(ctx, d, opts)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func (s *Service) Replace(ctx context.Context, key string, data interface{}) error {
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"_id": key}

	if _, err := s.collection.ReplaceOne(ctx, filter, data, opts); err != nil {
		return fmt.Errorf("failed to update document for key %s, error: %w", key, err)
	}
	return nil
}
