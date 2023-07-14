package shard

import "log"

type Shard struct {
	ShardCount int
	ShardID    int
	Shards     map[int][]string
}

func NewShard(shardCount int, socketAddr string, views []string) *Shard {
	shards := make(map[int][]string)
	s := Shard{
		ShardCount: shardCount,
		ShardID:    -1,
		Shards:     shards,
	}

	if len(views)/shardCount >= 2 {
		nodesInShard := len(views) / shardCount
		nodesSoFar := 0
		shardIdx := 1

		//Initialize empty list (of nodes) for each shard
		for i := 0; i < shardCount; i++ {
			s.Shards[i+1] = []string{}
		}

		for _, view := range views {
			if shardIdx <= shardCount {
				if view == socketAddr {
					s.ShardID = shardIdx
				}

				if nodesSoFar < nodesInShard {
					s.Shards[shardIdx] = append(s.Shards[shardIdx], view)
					nodesSoFar++
				} else {
					shardIdx++
					if shardIdx <= shardCount {
						nodesSoFar = 0
						s.Shards[shardIdx] = append(s.Shards[shardIdx], view)
						nodesSoFar++
						if view == socketAddr {
							s.ShardID = shardIdx
						}
					}
				}
			}
		}

		if (len(views) % shardCount) == 1 {
			shards[shardIdx-1] = append(shards[shardIdx-1], views[len(views)-1])
		}

	} else if s.ShardID == -1 {
		log.Print("shard count not specified")
	} else {
		panic("not enough nodes to have redundancy in shards. exiting program now")
	}

	return &s
}

func (s *Shard) Clear() {
	s.Shards = make(map[int][]string)
	s.ShardCount = 0
	s.ShardID = -1
}

func (s Shard) HashShardIndex(key string) int {
	return (sum([]byte(key)) % len(s.Shards)) + 1
}

func sum(bytes []byte) int {
	total := 0
	for _, b := range bytes {
		total += int(b)
	}

	return total
}
