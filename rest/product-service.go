package main

import (
	"database/sql"
	"errors"
	error2 "examples-go/checks/error"
	model2 "examples-go/models"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	host = "localhost"
	port = 5431
	u    = "postgres"
	p    = "postgres"
	db   = "postgres"
)

type ProductService interface {
	GetAll() []model2.Product
	GetByName(name string) model2.Product
	GetByID(id int) model2.Product
	Create(name string) int
	Delete(id int)
	Update(p model2.Product, id int) model2.Product
}

type Memory struct {
	m  map[int]model2.Product
	Id int
}

func InitMemoryStore() Memory {
	log.Println("Init memory store")
	m2 := make(map[int]model2.Product)
	return Memory{m2, 0}
}

func (m *Memory) GetAll() []model2.Product {
	products := make([]model2.Product, 0, len(m.m))
	for _, p := range m.m {
		products = append(products, p)
	}
	log.Println("Get all products =", products)
	return products
}

func (m *Memory) GetByName(name string) model2.Product {
	log.Println("Get product by name =", name)
	for _, p := range m.m {
		if p.Name == name {
			return p
		}
	}
	return model2.Product{}
}

func (m *Memory) GetById(id int) model2.Product {
	log.Println("Get product by id =", id)
	for _, p := range m.m {
		if p.Id == id {
			return p
		}
	}
	return model2.Product{}
}

func (m *Memory) Create(name string) int {
	id := 0
	v := m.m
	if v == nil {
		m.m = make(map[int]model2.Product)
	}
	id = m.nexId()
	group := model2.Group{}
	m.m[id] = model2.Product{Id: id, Name: name, Gr: group}
	log.Println("Create product by name =", name)
	return id
}

func (m *Memory) Delete(id int) {
	log.Println("Delete product by id=", id)
	delete(m.m, id)
}

func (m *Memory) Update(pr model2.Product, id int) model2.Product {
	product := m.m[id]
	product.Name = pr.Name
	product.Gr = pr.Gr
	m.m[id] = product
	log.Println("Update product by id =", id, ",product =", product)
	return product
}

func (m *Memory) nexId() int {
	if m.Id == 0 {
		m.Id = 1
	} else {
		m.Id++
	}
	return m.Id
}

func Create(name string) (string, error) {
	p := GetByName(name)
	if p.Id > 0 {
		msg := "Constrain product by name =" + name
		fmt.Println(msg)
		return "", errors.New(msg)
	}
	insertNewProduct := `insert into "product"("name") values($1)`
	connection := createConnection()
	defer connection.Close()
	_, e := connection.Exec(insertNewProduct, name)
	error2.CheckErrorDb(e)
	return name, nil
}

func GetByName(name string) model2.Product {
	if name != "" {
		c := createConnection()
		defer c.Close()
		product := model2.Product{}
		var groupId int
		row := c.QueryRow("SELECT id, name, coalesce(group_id, 0) FROM product WHERE name = $1", name)
		err := row.Scan(&product.Id, &product.Name, &groupId)
		error2.CheckErrorDb(err)
		product.Gr = findGroup(groupId)
		return product
	}
	return model2.Product{}
}

func GetById(id int) model2.Product {
	findProductByID := "SELECT id, name, coalesce(group_id, 0) FROM product WHERE id = $1"
	c := createConnection()
	defer c.Close()
	var name string
	var idDb int
	var groupId int
	err := c.QueryRow(findProductByID, id).Scan(&idDb, &name, &groupId)
	product := model2.Product{Id: idDb, Name: name, Gr: findGroup(groupId)}
	error2.CheckErrorDb(err)
	return product
}

func GetAll() []model2.Product {
	c := createConnection()
	defer c.Close()
	query, err := c.Query("SELECT id, name, COALESCE(group_id, 0) FROM product")
	error2.CheckErrorDb(err)
	products := make([]model2.Product, 0)
	for query.Next() {
		var name string
		var id int
		var groupId int
		err := query.Scan(&id, &name, &groupId)
		error2.CheckErrorDb(err)
		product := model2.Product{
			Id:   id,
			Name: name,
		}
		product.Gr = findGroup(id)
		products = append(products, product)
	}
	return products
}

func findGroup(id int) model2.Group {
	c := createConnection()
	defer c.Close()
	var group model2.Group
	if id > 0 {
		findGroupById := "SELECT * FROM product_group WHERE id = $1"
		err := c.QueryRow(findGroupById, id).Scan(&group.Id, &group.Name)
		error2.CheckErrorDb(err)
	}
	return group
}

func createConnection() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, u, p, db)
	c, err := sql.Open("postgres", psqlconn)
	error2.CheckErrorDb(err)
	_, err = c.Exec(`set search_path='test-go'`)
	error2.CheckErrorDb(err)
	return c
}
