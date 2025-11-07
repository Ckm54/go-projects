package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	tests := []struct {
		name      string
		ttl       time.Duration
		key       string
		val       []byte
		wait      time.Duration
		wantValue []byte
		wantOk    bool
	}{
		{
			name:      "basic add and get",
			ttl:       2 * time.Second,
			key:       "pokemon/1",
			val:       []byte("pikachu"),
			wait:      0,
			wantValue: []byte("pikachu"),
			wantOk:    true,
		},
		{
			name:      "expired entry",
			ttl:       100 * time.Millisecond,
			key:       "poke",
			val:       []byte("charmander"),
			wait:      200 * time.Millisecond,
			wantValue: nil,
			wantOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewCache(tt.ttl)
			cache.Add(tt.key, tt.val)

			if tt.wait > 0 {
				time.Sleep(tt.wait)
			}

			got, ok := cache.Get(tt.key)
			if ok != tt.wantOk {
				t.Errorf("cache.Get() ok = %v, want %v", ok, tt.wantOk)
			}
			if tt.wantOk && string(got) != string(tt.wantValue) {
				t.Errorf("cache.Get() = %s, want %s", got, tt.wantValue)
			}
		})
	}
}
