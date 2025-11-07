module go-redis-messaging-data-structures

go 1.24.0

toolchain go1.24.8

require (
	chatapp v0.0.0-00010101000000-000000000000
	github.com/redis/go-redis/v9 v9.14.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
)

replace chatapp => ./apps/chat-app
