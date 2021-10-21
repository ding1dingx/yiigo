package yiigo

import (
	"path/filepath"
	"strings"
	"sync"
)

type cfgdb struct {
	name    string
	driver  DBDriver
	dsn     string
	options []DBOption
}

type cfgmongo struct {
	name string
	dsn  string
}

type cfgredis struct {
	name    string
	address string
	options []RedisOption
}

type cfgnsq struct {
	nsqd    string
	lookupd []string
	options []NSQOption
}

type cfglogger struct {
	name    string
	path    string
	options []LoggerOption
}

type initSetting struct {
	logger []*cfglogger
	db     []*cfgdb
	mongo  []*cfgmongo
	redis  []*cfgredis
	nsq    *cfgnsq
}

// InitOption configures how we set up the yiigo initialization.
type InitOption func(s *initSetting)

// WithDB register db.
// [MySQL] username:password@tcp(localhost:3306)/dbname?timeout=10s&charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local
// [Postgres] host=localhost port=5432 user=root password=secret dbname=test connect_timeout=10 sslmode=disable
// [SQLite] file::memory:?cache=shared
func WithDB(name string, driver DBDriver, dsn string, options ...DBOption) InitOption {
	return func(s *initSetting) {
		s.db = append(s.db, &cfgdb{
			name:    name,
			driver:  driver,
			dsn:     dsn,
			options: options,
		})
	}
}

// WithMongo register mongodb.
// [DSN] mongodb://localhost:27017/?connectTimeoutMS=10000&minPoolSize=10&maxPoolSize=20&maxIdleTimeMS=60000&readPreference=primary
// [reference] https://docs.mongodb.com/manual/reference/connection-string
func WithMongo(name string, dsn string) InitOption {
	return func(s *initSetting) {
		s.mongo = append(s.mongo, &cfgmongo{
			name: name,
			dsn:  dsn,
		})
	}
}

// WithRedis register redis.
func WithRedis(name, address string, options ...RedisOption) InitOption {
	return func(s *initSetting) {
		s.redis = append(s.redis, &cfgredis{
			name:    name,
			address: address,
			options: options,
		})
	}
}

// WithNSQ initialize the nsq.
func WithNSQ(nsqd string, lookupd []string, options ...NSQOption) InitOption {
	return func(s *initSetting) {
		s.nsq = &cfgnsq{
			nsqd:    nsqd,
			lookupd: lookupd,
			options: options,
		}
	}
}

// WithLogger register logger.
func WithLogger(name, logfile string, options ...LoggerOption) InitOption {
	return func(s *initSetting) {
		cfg := &cfglogger{
			name:    name,
			options: options,
		}

		if v := strings.TrimSpace(logfile); len(v) != 0 {
			cfg.path = filepath.Clean(v)
		}

		s.logger = append(s.logger, cfg)
	}
}

// Init yiigo initialization.
func Init(options ...InitOption) {
	setting := new(initSetting)

	for _, f := range options {
		f(setting)
	}

	if len(setting.logger) != 0 {
		for _, v := range setting.logger {
			initLogger(v.name, v.path, v.options...)
		}
	}

	var wg sync.WaitGroup

	if len(setting.db) != 0 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for _, v := range setting.db {
				initDB(v.name, v.driver, v.dsn, v.options...)
			}
		}()
	}

	if len(setting.mongo) != 0 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for _, v := range setting.mongo {
				initMongoDB(v.name, v.dsn)
			}
		}()
	}

	if len(setting.redis) != 0 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for _, v := range setting.redis {
				initRedis(v.name, v.address, v.options...)
			}
		}()
	}

	if setting.nsq != nil {
		wg.Add(1)

		go func() {
			defer wg.Done()

			initNSQ(setting.nsq.nsqd, setting.nsq.lookupd, setting.nsq.options...)
		}()
	}

	wg.Wait()
}
