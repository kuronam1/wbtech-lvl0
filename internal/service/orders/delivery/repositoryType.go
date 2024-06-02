package delivery

type Repository interface {
	CreateDelivery()
	GetOne()
	GetAllDeliveries()
}
