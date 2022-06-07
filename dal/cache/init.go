package cache

func Init() {
	t := NewMyTick(300, redisToSQL)
	go t.Start()
}
