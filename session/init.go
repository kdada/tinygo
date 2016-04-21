package session

func init() {
	//注册 内存SessionContainer创建器
	Register("memory", newMemSessionContainer)
}
