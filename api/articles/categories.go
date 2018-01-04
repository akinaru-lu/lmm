package articles

import (
	"lmm/api/db"
	"net/http"

	"github.com/akinaru-lu/elesion"
)

type Category struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

func GetCategories(c *elesion.Context) {
	userID := c.Query().Get("user_id")
	if userID == "" {
		c.Status(http.StatusBadRequest).String("missing user_id")
		return
	}

	categories, err := getCategories(userID)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(categories)
}

func getCategories(userID string) ([]Category, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	itr, err := d.Query("SELECT id, user_id, name FROM categories WHERE user_id = ? ORDER BY name", userID)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	categories := make([]Category, 0)
	for itr.Next() {
		category := Category{}
		itr.Scan(&category.ID, &category.UserID, &category.Name)

		categories = append(categories, category)
	}
	return categories, nil
}

func GetCategoryByID(c *elesion.Context) {
	id := c.Query().Get("id")
	if id == "" {
		c.Status(http.StatusBadRequest).String("missing id")
		return
	}

	category, err := getCategoryByID(id)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(category)
}

func getCategoryByID(id string) (*Category, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	category := new(Category)
	err := d.QueryRow("SELECT id, user_id, name FROM categories WHERE id = ?", id).Scan(
		&category.ID, &category.UserID, &category.Name,
	)
	return category, err
}

func getCategoryIDByName(name string) (string, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	id := ""
	err := d.QueryRow("SELECT id FROM categories WHERE name = ?", name).Scan(&id)

	return id, err
}
