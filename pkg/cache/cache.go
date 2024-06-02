package cache

type Caсher interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
}

func InitCache() Caсher {
	items := make(map[string]Item)

	c := &Cache{
		Items: items,
	}

	return c
}

func (c *Cache) Set(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()

	c.Items[key] = value
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	value, found := c.Items[key]
	return value, found
}

/*func (c *Cache) Exists(key string) bool {
	_, found := c.Items[key]
	if !found {
		return false
	} else {
		return true
	}
}*/

/*func (c *Cache) Delete(key string) error {
	return nil
}

func (c *Cache) StartGarbageCleaner(ctx context.Context) {
	go c.GarbageCleaner(ctx)
}

func (c *Cache) GarbageCleaner(ctx context.Context) {
	for {
		select {
		case <-time.After(c.CleanupInterval):
			if c.Items == nil {
				c.logger.Error("[cache] Garbage Cleaner - Items == nil")
				return
			}

			if keys := c.ExpiredKeys(); len(keys) > 0 {
				c.ClearItems(keys)
			}

		case <-ctx.Done():
			c.logger.Info("[cache] Stopping Garbage Cleaner...")
			return
		}
	}
}

func (c *Cache) ExpiredKeys() (keys []string) {
	c.RLock()
	defer c.RUnlock()

	for key, value := range c.Items {
		if time.Now().Unix() > value.Expire && value.Expire > 0 {
			keys = append(keys, key)
		}
	}
	return
}

func (c *Cache) ClearItems(keys []string) {
	c.Lock()
	defer c.Unlock()

	for _, item := range keys {
		delete(c.Items, item)
	}
}*/
