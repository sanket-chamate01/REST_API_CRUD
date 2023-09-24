package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type products struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float32 `json:"price"`
}

func getProducts(db *sql.DB) ([]products, error) {
	query := "select id, name, quantity, price from products"
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	product := []products{}
	for rows.Next() {
		var p products
		err := rows.Scan(&p.Id, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		product = append(product, p)
	}
	return product, nil
}

func (p *products) getProductId(db *sql.DB) error {
	query := fmt.Sprintf("select name, quantity, price from products where id=%v", p.Id)
	row := db.QueryRow(query)
	err := row.Scan(&p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	return nil
}

func (p *products) createProduct(db *sql.DB) error {
	query := fmt.Sprintf("insert into products(name, quantity, price) values('%v',%v,%v)", p.Name, p.Quantity, p.Price)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.Id = int(id)
	return nil
}

func (p *products) updateProduct(db *sql.DB) error {
	query := fmt.Sprintf("update products set name='%v', quantity=%v, price=%v where id=%v", p.Name, p.Quantity, p.Price, p.Id)
	result, err := db.Exec(query)
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("No such row exists")
	}
	return err
}

func (p *products) deleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("delete from products where id=%v", p.Id)
	result, err := db.Exec(query)
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("No such row exists")
	}
	return err
}
