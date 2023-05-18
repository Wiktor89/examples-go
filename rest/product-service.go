package main

import (
	"database/sql"
	"errors"
	er "examples-go/checks/error"
	model "examples-go/models"
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
	GetAll() []model.Product
	GetByName(name string) model.Product
	GetByID(id int) model.Product
	Create(name string) int
	Delete(id int)
	Update(p model.Product, id int) model.Product
}

type Memory struct {
	m  map[int]model.Product
	Id int
}

func InitMemoryStore() Memory {
	log.Println("Init memory store")
	m2 := make(map[int]model.Product)
	return Memory{m2, 0}
}

func (m *Memory) GetAll() []model.Product {
	products := make([]model.Product, 0, len(m.m))
	for _, p := range m.m {
		products = append(products, p)
	}
	log.Println("Get all products =", products)
	return products
}

func (m *Memory) GetByName(name string) model.Product {
	log.Println("Get product by name =", name)
	for _, p := range m.m {
		if p.Name == name {
			return p
		}
	}
	return model.Product{}
}

func (m *Memory) GetById(id int) model.Product {
	log.Println("Get product by id =", id)
	for _, p := range m.m {
		if p.Id == id {
			return p
		}
	}
	return model.Product{}
}

func (m *Memory) Create(name string) int {
	id := 0
	v := m.m
	if v == nil {
		m.m = make(map[int]model.Product)
	}
	id = m.nexId()
	group := model.Group{}
	m.m[id] = model.Product{Id: id, Name: name, Gr: group}
	log.Println("Create product by name =", name)
	return id
}

func (m *Memory) Delete(id int) {
	log.Println("Delete product by id=", id)
	delete(m.m, id)
}

func (m *Memory) Update(pr model.Product, id int) model.Product {
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
	er.CheckErrorDb(e)
	return name, nil
}

func GetByName(name string) model.Product {
	if name != "" {
		c := createConnection()
		defer c.Close()
		product := model.Product{}
		var groupId int
		row := c.QueryRow("SELECT id, name, coalesce(group_id, 0) FROM product WHERE name = $1", name)
		err := row.Scan(&product.Id, &product.Name, &groupId)
		er.CheckErrorDb(err)
		product.Gr = findGroup(groupId)
		return product
	}
	return model.Product{}
}

func GetById(id int) model.Product {
	findProductByID := "SELECT id, name, coalesce(group_id, 0) FROM product WHERE id = $1"
	c := createConnection()
	defer c.Close()
	var name string
	var idDb int
	var groupId int
	err := c.QueryRow(findProductByID, id).Scan(&idDb, &name, &groupId)
	product := model.Product{Id: idDb, Name: name, Gr: findGroup(groupId)}
	er.CheckErrorDb(err)
	return product
}

func GetAll() []model.Product {
	c := createConnection()
	defer c.Close()
	query, err := c.Query("SELECT id, name, COALESCE(group_id, 0) FROM product")
	er.CheckErrorDb(err)
	products := make([]model.Product, 0)
	for query.Next() {
		var name string
		var id int
		var groupId int
		err := query.Scan(&id, &name, &groupId)
		er.CheckErrorDb(err)
		product := model.Product{
			Id:   id,
			Name: name,
		}
		product.Gr = findGroup(id)
		products = append(products, product)
	}
	return products
}

func findGroup(id int) model.Group {
	c := createConnection()
	defer c.Close()
	var group model.Group
	if id > 0 {
		findGroupById := "SELECT * FROM product_group WHERE id = $1"
		err := c.QueryRow(findGroupById, id).Scan(&group.Id, &group.Name)
		er.CheckErrorDb(err)
	}
	return group
}

func createConnection() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, u, p, db)
	c, err := sql.Open("postgres", psqlconn)
	er.CheckErrorDb(err)
	_, err = c.Exec(`set search_path='test-go'`)
	er.CheckErrorDb(err)
	return c
}
