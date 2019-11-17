package redisbloom

import (
	redis "github.com/go-redis/redis/v7"
)



func ExtendClient(client *redis.Client) *Client {
	return &Client {
		client,
		&redisCmnder{Process:client.Process},
	}
}



type Client struct {
	*redis.Client
	*redisCmnder
}

type RbCmdable interface {
	BFReserve(key string, errorRate float64, capacity int64) *redis.StatusCmd
	BFAdd(key string, item interface{}) *redis.IntCmd
	BFMadd(key string, items ...interface{}) *redis.IntSliceCmd
	BFExists(key string, item interface{}) *redis.IntCmd
	BFMexists(key string, items ...interface{}) *redis.IntSliceCmd
	//TODO: BF.INSERT {key} [CAPACITY {cap}] [ERROR {error}] [NOCREATE] ITEMS {item...}
	//TODO: BF.SCANDUMP {key} {iter}
	//TODO: BF.LOADCHUNK {key} {iter} {data}
}
type redisCmnder struct {
	Process func(cmd redis.Cmder) error
}


func appendArgs(dst, src []interface{}) []interface{} {
	if len(src) == 1 {
		if ss, ok := src[0].([]string); ok {
			for _, s := range ss {
				dst = append(dst, s)
			}
			return dst
		}
	}

	dst = append(dst, src...)
	return dst
}


func(c *redisCmnder) BFReserve(key string,  errorRate float64, capacity int64) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(
		"bf.reserve",
		key,
		errorRate,
		capacity,
	)
	_ = c.Process(cmd)
	return cmd

}

func(c redisCmnder) BFAdd(key string, item interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd("bf.add", key, item)
	_ = c.Process(cmd)
	return cmd
}

func(c redisCmnder) BFMadd(key string, items ...interface{}) *redis.IntSliceCmd {
	args := make([]interface{}, 2, 2+len(items))
	args[0] = "bf.add"
	args[1] = key
	args = appendArgs(args, items)
	cmd := redis.NewIntSliceCmd(args...)
	_ = c.Process(cmd)
	return cmd
}

func(c redisCmnder) BFExists(key string, item interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd("bf.exists", key, item)
	_ = c.Process(cmd)
	return cmd
}
func(c redisCmnder) BFMexists(key string, items ...interface{}) *redis.IntSliceCmd {
	args := make([]interface{}, 2, 2+len(items))
	args[0] = "bf.exists"
	args[1] = key
	args = appendArgs(args, items)
	cmd := redis.NewIntSliceCmd(args...)
	_ = c.Process(cmd)
	return cmd
}