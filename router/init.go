package router

func init() {
	Register("unlimited", NewUnlimitedRouter)
	Register("base", NewBaseRouter)
}
