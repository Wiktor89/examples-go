package main

import (
	"examples-go/models"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemory_Create(t *testing.T) {
	store := testInitMemory()
	id := store.Create("name")
	product, _ := store.GetById(id)
	assert.Equal(t, 1, id, "Id product not match.")
	assert.Equal(t, "name", product.Name, "Name product not match.")
}

func TestMemory_GetAll(t *testing.T) {
	store := testInitMemory()
	id := store.Create("name")
	assert.Equal(t, 1, id, "Id product not match.")
}

func TestMemory_GetByName(t *testing.T) {
	store := testInitMemory()
	store.Create("name")
	product := store.GetByName("name")
	assert.Equal(t, "name", product.Name, "Name product not match.")
}

func TestMemory_GetById(t *testing.T) {
	store := testInitMemory()
	store.Create("name")
	product, _ := store.GetById(1)
	assert.Equal(t, 1, product.Id, "Get by id product.")
}

func TestMemory_Delete(t *testing.T) {
	store := testInitMemory()
	idProduct := store.Create("name")
	store.Delete(idProduct)
	product, _ := store.GetById(idProduct)
	assert.Equal(t, models.Product{}, product, "Product not exist.")
}

func TestMemory_Update(t *testing.T) {
	store := testInitMemory()
	idProduct := store.Create("name")
	product, _ := store.GetById(idProduct)
	assert.Equal(t, product.Id, 1, fmt.Sprintf("Not product by id = %d", idProduct))
	assert.Equal(t, product.Name, "name", "Product not exist name.")
	store.Update(models.Product{Name: "change name"}, product.Id)
	productUpdate, _ := store.GetById(idProduct)
	assert.Equal(t, "change name", productUpdate.Name, "Product name not change after update.")
}

func testInitMemory() Memory {
	return InitMemoryStore()
}
