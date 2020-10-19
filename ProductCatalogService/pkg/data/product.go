package product

//Product data model
type Product struct {
	ID          string
	Name        string
	Description string
	Picture     string
	Price       float32
	Categories  []string
}
