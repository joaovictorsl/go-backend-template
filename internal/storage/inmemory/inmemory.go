package inmemory

import "sync"

type KVCache struct {
	data *sync.Map
}

func New() *KVCache {
	return &KVCache{
		data: &sync.Map{},
	}
}

func (c *KVCache) Get(key string) (string, error) {
	v, ok := c.data.Load(key)
	if !ok {
		return "", nil
	}
	return v.(string), nil
}

func (c *KVCache) Insert(key string, value string) error {
	c.data.Store(key, value)
	return nil
}

func (c *KVCache) Remove(key string) {
	c.data.Delete(key)
}
