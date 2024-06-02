package items

type Repository interface {
	CreateItem()
	GetOne()
	GetAllItemsByOrderID()
}
