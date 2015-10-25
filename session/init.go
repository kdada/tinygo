package session

func init() {
	//注册内存会话提供器
	RegisterSessionProviderCreator(SessionTypeMemory, newMemSessionProvider)
}
