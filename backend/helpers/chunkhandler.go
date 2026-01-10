package helpers

import "sync"

type ChunkEventTracker struct {
	m    sync.Mutex
	data map[string]map[string]any
}

var ChunkData = ChunkEventTracker{
	data: map[string]map[string]any{},
}

func SaveChunkData(id string, data map[string]any) {
	ChunkData.m.Lock()
	ChunkData.data[id] = data
	ChunkData.m.Unlock()
}

func ReadChunkData(id string) map[string]any {
	ChunkData.m.Lock()
	defer ChunkData.m.Unlock()
	return ChunkData.data[id]
}
