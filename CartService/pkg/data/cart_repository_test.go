package data

import (
	"testing"

	"github.com/ajanthan/go-ecommerce-demo/CartService/pkg/model"
	_ "github.com/go-sql-driver/mysql"
)

var cartRepository = &CartRepository{}

const dsn = "root:root@tcp(localhost:3306)/cartdb"

func TestCartRepository_InitRepository(t *testing.T) {
	cartRepository.InitRepository(dsn)
}

func TestCartRepository_AddItemToCart(t *testing.T) {
	item1 := model.Item{}
	item1.ProductID = "fh58tt5"
	item1.Quantity = 1
	cartRepository.AddItemToCart("testUser1", item1)
	item2 := model.Item{}
	item2.ProductID = "fh58tt6"
	item2.Quantity = 2
	cartRepository.AddItemToCart("testUser1", item2)
	item3 := model.Item{}
	item3.ProductID = "fh58tt7"
	item3.Quantity = 3
	cartRepository.AddItemToCart("testUser1", item3)
	item4 := model.Item{}
	item4.ProductID = "fh58tt5"
	item4.Quantity = 2
	cartRepository.AddItemToCart("testUser2", item4)
	item5 := model.Item{}
	item5.ProductID = "fh58tt6"
	item5.Quantity = 3
	cartRepository.AddItemToCart("testUser2", item5)

}

func TestCartRepository_GetCart(t *testing.T) {
	cart1 := cartRepository.GetCart("testUser1")
	if len(cart1.Items) != 3 {
		t.FailNow()
	}
	cart2 := cartRepository.GetCart("testUser2")
	if len(cart2.Items) != 2 {
		t.FailNow()
	}
}

func TestCartRepository_EmptyCart(t *testing.T) {
	cartRepository.EmptyCart("testUser1")
	cart1 := cartRepository.GetCart("testUser1")
	if len(cart1.Items) != 0 {
		t.FailNow()
	}
}
