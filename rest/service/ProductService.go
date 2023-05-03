package service

import (
	"database/sql"
	"examples-go/common/model"
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

func Create(name string) {
	insertNewProduct := `insert into "product"("name") values($1)`
	connection := createConnection()
	defer connection.Close()
	_, e := connection.Exec(insertNewProduct, name)
	CheckError(e)
}

func GetById(id int) model.Product {
	findProductByID := "SELECT id, name, coalesce(group_id, 0) FROM product WHERE id = $1"
	connection := createConnection()
	defer connection.Close()
	var name string
	var idDb int
	var groupId int
	err := connection.QueryRow(findProductByID, id).Scan(&idDb, &name, &groupId)
	product := model.Product{Id: idDb, Name: name, Gr: findGroup(groupId, connection)}
	CheckError(err)
	return product
}

func CheckError(err error) {
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			log.Println("No rows were returned!")
			break
		case nil:
			break
		default:
			log.Fatalln("Unable to scan the row. %v", err)
		}
	}
}

func GetAll() []model.Product {
	connection := createConnection()
	defer connection.Close()
	query, err := connection.Query("SELECT id, name, COALESCE(group_id, 0) FROM product")
	CheckError(err)
	products := make([]model.Product, 0)
	for query.Next() {
		var name string
		var idDb int
		var groupId int
		err := query.Scan(&idDb, &name, &groupId)
		CheckError(err)
		products = append(products, model.Product{
			Id:   idDb,
			Name: name,
			Gr:   findGroup(groupId, connection),
		})
	}
	return products
}

func findGroup(groupId int, c *sql.DB) model.Group {
	var group model.Group
	if groupId > 0 {
		var grId int
		var grName string
		findGroupById := "SELECT * FROM product_group WHERE id = $1"
		c.QueryRow(findGroupById, groupId).Scan(&grId, &grName)
		group.Id = grId
		group.Name = grName
	}
	return group
}

func createConnection() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, u, p, db)
	open, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	open.Exec(`set search_path='test-go'`)
	return open
}
