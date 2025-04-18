package service

type PaymentService interface {
	CreateLink()
	ApprovePayment()
	CancelPayment()
}

type UserService interface {
	Create()
	User()
}

type SubscriptionService interface {
}
