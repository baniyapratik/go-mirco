package mongo_driver

import "log"

type Option func(s *Service)

func WithLogger(logger log.Logger) Option {
	return func(s *Service) {
		s.logger = logger
	}
}

func WithMongoDBURL(url string) Option {
	return func(s *Service) {
		s.dbURL = url
	}
}

func WithUsername(u string) Option {
	return func(s *Service) {
		s.username = u
	}
}

func WithPassword(p string) Option {
	return func(s *Service) {
		s.password = p
	}
}

func WithDatabaseName(d string) Option {
	return func(s *Service) {
		s.databaseName = d
	}
}

func WithCollectionName(c string) Option {
	return func(s *Service) {
		s.collectionName = c
	}
}

func WithCACertFile(ca string) Option {
	return func(s *Service) {
		s.caCertFile = ca
	}
}
