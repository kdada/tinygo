package connector

func init() {
	Register("http", NewHttpConnector)
}
