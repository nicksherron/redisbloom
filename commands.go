package redisbloom

import (
	redis "github.com/go-redis/redis/v7"
)

type Client struct {
	*redis.Client
	*cmdable
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

type cmdable struct {
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


func(c cmdable) BFReserve(key string,  errorRate float64, capacity int64) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(
		"bf.reserve",
		key,
		errorRate,
		capacity,
	)
	_ = c(cmd)
	return cmd

}

func(c cmdable) BFAdd(key string, item interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd("bf.add", key, item)
	_ = c(cmd)
	return cmd
}

func(c cmdable) BFMadd(key string, items ...interface{}) *redis.IntSliceCmd {
	args := make([]interface{}, 2, 2+len(items))
	args[0] = "bf.add"
	args[1] = key
	args = appendArgs(args, items)
	cmd := redis.NewIntSliceCmd(args...)
	_ = c(cmd)
	return cmd
}

func(c cmdable) BFExists(key string, item interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd("bf.exists", key, item)
	_ = c(cmd)
	return cmd
}
func(c cmdable) BFMexists(key string, items ...interface{}) *redis.IntSliceCmd {
	args := make([]interface{}, 2, 2+len(items))
	args[0] = "bf.exists"
	args[1] = key
	args = appendArgs(args, items)
	cmd := redis.NewIntSliceCmd(args...)
	_ = c(cmd)
	return cmd
}