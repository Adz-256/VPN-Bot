package text

const (
	Subscriptions = `📋 Ваши подписки:`
)

func Show(id string) string {
	return `📋 Подписка: ` + id
}
