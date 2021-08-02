package yiigo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/shenghui0779/vitess_pool"
	"go.uber.org/zap"
)

type redisConfig struct {
	Address      string `toml:"address"`
	Password     string `toml:"password"`
	Database     int    `toml:"database"`
	ConnTimeout  int    `toml:"conn_timeout"`
	ReadTimeout  int    `toml:"read_timeout"`
	WriteTimeout int    `toml:"write_timeout"`
	PoolSize     int    `toml:"pool_size"`
	PoolLimit    int    `toml:"pool_limit"`
	IdleTimeout  int    `toml:"idle_timeout"`
	PoolPrefill  int    `toml:"pool_prefill"`
}

// RedisConn redis connection resource
type RedisConn struct {
	redis.Conn
}

// Close close connection resorce
func (rc *RedisConn) Close() {
	if err := rc.Conn.Close(); err != nil {
		logger.Error("yiigo: redis conn closed error", zap.Error(err))
	}
}

// RedisPool redis pool resource
type RedisPool interface {
	// Get get a connection resource from the pool.
	// Context with timeout can specify the wait timeout for pool.
	Get(ctx context.Context) (*RedisConn, error)

	// Put returns a connection resource to the pool.
	Put(rc *RedisConn)
}

type redisPoolResource struct {
	config *redisConfig
	pool   *vitess_pool.ResourcePool
	mutex  sync.Mutex
}

func (r *redisPoolResource) dial() (redis.Conn, error) {
	dialOptions := []redis.DialOption{
		redis.DialPassword(r.config.Password),
		redis.DialDatabase(r.config.Database),
		redis.DialConnectTimeout(time.Duration(r.config.ConnTimeout) * time.Second),
		redis.DialReadTimeout(time.Duration(r.config.ReadTimeout) * time.Second),
		redis.DialWriteTimeout(time.Duration(r.config.WriteTimeout) * time.Second),
	}

	conn, err := redis.Dial("tcp", r.config.Address, dialOptions...)

	return conn, err
}

func (r *redisPoolResource) init() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.pool != nil && !r.pool.IsClosed() {
		return
	}

	df := func() (vitess_pool.Resource, error) {
		conn, err := r.dial()

		if err != nil {
			return nil, err
		}

		return &RedisConn{conn}, nil
	}

	r.pool = vitess_pool.NewResourcePool(df, r.config.PoolSize, r.config.PoolLimit, time.Duration(r.config.IdleTimeout)*time.Second, r.config.PoolPrefill)
}

func (r *redisPoolResource) Get(ctx context.Context) (*RedisConn, error) {
	if r.pool.IsClosed() {
		r.init()
	}

	resource, err := r.pool.Get(ctx)

	if err != nil {
		return &RedisConn{}, err
	}

	rc := resource.(*RedisConn)

	// If rc is error, close and reconnect
	if rc.Err() != nil {
		conn, err := r.dial()

		if err != nil {
			r.pool.Put(rc)

			return rc, err
		}

		rc.Close()

		return &RedisConn{conn}, nil
	}

	return rc, nil
}

func (r *redisPoolResource) Put(rc *RedisConn) {
	r.pool.Put(rc)
}

var (
	defaultRedis RedisPool
	redisMap     sync.Map
)

func initRedis() {
	configs := make(map[string]*redisConfig)

	if err := env.Get("redis").Unmarshal(&configs); err != nil {
		logger.Panic("yiigo: redis init error", zap.Error(err))
	}

	if len(configs) == 0 {
		return
	}

	for name, cfg := range configs {
		pool := &redisPoolResource{config: cfg}

		pool.init()

		// verify connection
		conn, err := pool.Get(context.TODO())

		if err != nil {
			logger.Panic("yiigo: redis init error", zap.String("name", name), zap.Error(err))
		}

		if _, err = conn.Do("PING"); err != nil {
			conn.Close()

			logger.Panic("yiigo: redis init error", zap.String("name", name), zap.Error(err))
		}

		pool.Put(conn)

		if name == defaultConn {
			defaultRedis = pool
		}

		redisMap.Store(name, pool)

		logger.Info(fmt.Sprintf("yiigo: redis.%s is OK.", name))
	}
}

// Redis returns a redis pool.
func Redis(name ...string) RedisPool {
	if len(name) == 0 {
		if defaultRedis == nil {
			logger.Panic(fmt.Sprintf("yiigo: unknown redis.%s (forgotten configure?)", defaultConn))
		}

		return defaultRedis
	}

	v, ok := redisMap.Load(name[0])

	if !ok {
		logger.Panic(fmt.Sprintf("yiigo: unknown redis.%s (forgotten configure?)", name[0]))
	}

	return v.(RedisPool)
}
