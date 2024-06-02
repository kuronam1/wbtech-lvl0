package payment

type Repository interface {
	CreatePayment()
	GetOnePayment()
	GetAllPayments()
}
