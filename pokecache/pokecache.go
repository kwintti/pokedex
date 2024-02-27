package pokecache

import(
    "time"
    "sync"
)

func NewCache(interval time.Duration) *Cache {
    c := &Cache{
                    m: make(map[string]cacheEntry),
                    mu: sync.Mutex{},
                    interval: interval * time.Second,
                }
    go c.reapLoop()
    return c 
}

func (c *Cache) Add(key string, val []byte) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.m[key] = cacheEntry{createdAt: time.Now(), val: val}   
}

func (c *Cache) Get(key string) ([]byte, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if v, ok := c.m[key]; ok {
        return v.val, true
    }
    return nil, false
}

func (c *Cache) reapLoop() {
    ticker := time.NewTicker(c.interval)
    for t := range ticker.C {
        c.mu.Lock()
        for k, item := range c.m {
            if t.After(item.createdAt) {
                delete(c.m, k)
            }
        }
        c.mu.Unlock()
    }
}

type Cache struct {
    m           map[string]cacheEntry
    mu          sync.Mutex
    interval    time.Duration
}

type cacheEntry struct {
    createdAt   time.Time
    val         []byte
}
