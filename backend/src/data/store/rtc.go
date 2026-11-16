package store

import "zzc/fall-chat/src/data/cache"

type RTCStore interface {
	RTCCache() cache.RTCCache
}
