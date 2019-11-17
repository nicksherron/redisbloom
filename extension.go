package redisbloom

import "github.com/go-redis/redis"

func ExtendClient(client *redis.Client) *Client {
	return &Client {
		client,
		&cmdable {
			Process: client.Process,
		},
	}
}
