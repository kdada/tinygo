package connector

func init() {
	Register("http", NewHttpConnector)
	Register("https", NewHttpsConnector)
}
