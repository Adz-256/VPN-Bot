package webhook

type Webhook interface {
	Run() chan map[string]any
}
