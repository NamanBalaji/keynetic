package shard

type Shard struct {
	shardCount int
	shardID    int
	shards     map[int][]string
}

func NewShard(shardCount int, socketAddr string, views []string) *Shard {
	shards := make(map[int][]string)
	s := Shard{
		shardCount: shardCount,
		shardID:    -1,
		shards:     shards,
	}

	if len(views)/shardCount >= 2 {
		nodesInShard := len(views) / shardCount
		nodesSoFar := 0
		shardIdx := 1

		//Initialize empty list (of nodes) for each shard
		for i := 0; i < shardCount; i++ {
			s.shards[i+1] = []string{}
		}

		for _, view := range views {
			if shardIdx <= shardCount {
				if view == socketAddr {
					s.shardID = shardIdx
				}

				if nodesSoFar < nodesInShard {
					s.shards[shardIdx] = append(s.shards[shardIdx], view)
					nodesSoFar++
				} else {
					shardIdx++
					if shardIdx <= shardCount {
						nodesSoFar = 0
						s.shards[shardIdx] = append(s.shards[shardIdx], view)
						nodesSoFar++
						if view == socketAddr {
							s.shardID = shardIdx
						}
					}
				}
			}
		}

		if (len(views) % shardCount) == 1 {
			shards[shardIdx-1] = append(shards[shardIdx-1], views[len(views)-1])
		}

	} else if s.shardID == -1 {
		panic("shard count not specified")
	} else {
		panic("not enough nodes to have redundancy in shards. exiting program now")
	}

	return &s
}
