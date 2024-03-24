package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Stock int `json:"stock"`
}


var products = make(map[string]Product)

func getAllProducts(c echo.Context) error {
	productsList := make([]Product, 0, len(products))
	for _, v := range products {
		productsList = append(productsList, v)
	}

	return c.JSON(http.StatusOK, productsList)
}

func createProduct(c echo.Context) error {
	newProduct := new(Product)
	if err := c.Bind(newProduct); err != nil {
		return err
	}

	if _, found := products[newProduct.ID]; found {
		return c.String(http.StatusBadRequest, "Product already exists")
	}

	products[newProduct.ID] = *newProduct
	return c.JSON(http.StatusCreated, newProduct)
}

func getProduct(c echo.Context) error {
	id := c.Param("id")
	product, found := products[id]
	if !found {
		return c.String(http.StatusNotFound, "Product not found")
	}

	return c.JSON(http.StatusOK, product)
}

func updateProduct(c echo.Context) error {
	id := c.Param("id")
	if _, found := products[id]; !found {
		return c.String(http.StatusNotFound, "Product not found")
	}

	updatedProduct := new(Product)
	if err := c.Bind(updatedProduct); err != nil {
		return err
	}

	products[id] = *updatedProduct
	return c.NoContent(http.StatusOK)
}

func deleteProduct(c echo.Context) error {
	id := c.Param("id")
	if _, found := products[id]; !found {
		return c.String(http.StatusNotFound, "Product not found")
	}

	delete(products, id)
	return c.NoContent(http.StatusOK)
}

func main() {
	e := echo.New()

	e.GET("/products", getAllProducts)
	e.POST("/product", createProduct)
	e.GET("/product/:id", getProduct)
	e.PUT("/product/:id", updateProduct)
	e.DELETE("/product/:id", deleteProduct)

	e.Logger.Fatal(e.Start("localhost:1234"))
}