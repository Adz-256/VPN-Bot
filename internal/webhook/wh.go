package webhook

type Webhook interface {
	Run()
	Channel() <-chan map[string]interface{}
}
