package pkg

import (
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"time"
)

type Client struct {
	cli    *redis.Client
	logger *zap.Logger
}

func NewClient(addr, password string, db int, logger *zap.Logger) (*Client, error) {
	cli := redis.NewClient(&redis.Options{
		Network:     "tcp",
		Addr:        addr,
		Password:    password,
		DB:          db,
		PoolSize:    10,
		MaxRetries:  2,
		IdleTimeout: 10 * time.Minute,
	})
	
	return &Client{
		cli:    cli,
		logger: logger,
	}, cli.Ping().Err()
}

// 删除key

func (client *Client) Del(key string) {
	//
	t, err := client.cli.Type(key).Result()
	if err != nil {
		client.logger.Error(
			"获取redis key类型失败",
			zap.Error(err),
		)
		return
	}
	
	// key not found 时，type 为 none
	if t != "string" {
		client.logger.Error(
			"仅能获取string类型的key",
			zap.String("type", t),
		)
		return
	}
	
	//
	err = client.cli.Del(key).Err()
	if err != nil {
		client.logger.Error(
			"删除redis key失败",
			zap.Error(err),
		)
		return
	}
	
	fmt.Printf("key: %s, type: %s 删除成功.\n", key, t)
}

// 获取key

func (client *Client) Get(key string) {
	//
	t, err := client.cli.Type(key).Result()
	if err != nil {
		client.logger.Error(
			"获取redis key类型失败",
			zap.Error(err),
		)
		return
	}
	
	// key not found 时，type 为 none
	if t != "string" {
		client.logger.Error(
			"仅能获取string类型的key",
			zap.String("type", t),
		)
		return
	}
	
	//
	v, err := client.cli.Get(key).Result()
	if err != nil {
		client.logger.Error(
			"获取redis key失败",
			zap.Error(err),
		)
		return
	}
	
	fmt.Printf("key: %s, value: %s, type: %s\n", key, v, t)
}

// 扫描

func (client *Client) Scan(match string, count int64, keyType bool) {
	var cursor uint64 = 0
	for {
		keys, newCursor, err := client.cli.Scan(cursor, match, count).Result()
		
		if err != nil {
			break
		}
		
		for _, key := range keys {
			
			if keyType {
				keyType, err := client.cli.Type(key).Result()
				if err != nil {
					panic(fmt.Sprintf("err: %s", err))
				}
				fmt.Printf("key: %s, type: %s\n", key, keyType)
			} else {
				fmt.Printf("key: %s\n", key)
			}
		}
		
		if newCursor == 0 {
			break
		}
		
		cursor = newCursor
	}
}

func (client *Client) KeyScan(key, match string, count int64) {
	var cursor uint64 = 0
	
	//
	t, err := client.cli.Type(key).Result()
	if err != nil {
		client.logger.Error(
			"获取redis key类型失败",
			zap.Error(err),
		)
		return
	}
	
	//
	switch t {
	// hset / hgetall
	case "hash":
		for {
			keys, newCursor, err := client.cli.HScan(key, cursor, match, count).Result()
			if err != nil {
				client.logger.Error(
					"扫描 hash key 失败",
					zap.Error(err),
				)
				break
			}
			
			for i := 0; i < len(keys)-1; i += 2 {
				fmt.Printf("hash key: %s, value: %s\n", keys[i], keys[i+1])
			}
			
			if newCursor == 0 {
				break
			}
			
			cursor = newCursor
		}
	
	// sadd / smembers
	case "set":
		for {
			keys, newCursor, err := client.cli.SScan(key, cursor, match, count).Result()
			if err != nil {
				client.logger.Error(
					"扫描 set key 失败",
					zap.Error(err),
				)
				break
			}
			
			for _, key := range keys {
				fmt.Printf("set key: %s\n", key)
			}
			
			if newCursor == 0 {
				break
			}
			
			cursor = newCursor
		}
	// zadd / zrange
	case "zset":
		for {
			keys, newCursor, err := client.cli.ZScan(key, cursor, match, count).Result()
			if err != nil {
				client.logger.Error(
					"扫描 zset key 失败",
					zap.Error(err),
				)
			}
			
			for i := 0; i < len(keys)-1; i += 2 {
				fmt.Printf("zset key: %s, value: %s\n", keys[i], keys[i+1])
			}
			
			if newCursor == 0 {
				break
			}
			
			cursor = newCursor
		}
	}
}

func (client *Client) Info() {
	r, err := client.cli.Info().Result()
	if err != nil {
		client.logger.Error(
			"获取redis信息失败",
			zap.Error(err),
		)
		return
	}
	fmt.Println(r)
}

func (client *Client) ClientList() {
	r, err := client.cli.ClientList().Result()
	if err != nil {
		client.logger.Error(
			"获取redis客户端列表失败",
			zap.Error(err),
		)
		return
	}
	fmt.Println(r)
}
