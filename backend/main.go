package main

import (
	"errors"
	"fmt"
	"log"
	"net/http" // connect GIN
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin" // connect GIN
	"gorm.io/driver/mysql"     // connect GORM
	"gorm.io/gorm"             // connect GORM
)

// Use enum data types for better data binding (Status  String `json:"status" gorm:"column:status;"`)
type ItemStatus int

const (
	ItemStatusInStock ItemStatus = iota
	ItemStatusOutOfStock
)

var allItemStatus = [2]string{"In stock", "Out of stock"}

// Convert ItemStatus to string
func (item ItemStatus) String() string {
	return allItemStatus[item]
}

// Convert string back to ItemStatus
func parseStr2ItemStatus(s string) (ItemStatus, error) {
	for i := range allItemStatus {
		if allItemStatus[i] == s {
			return ItemStatus(i), nil
		}
	}

	return ItemStatus(0), errors.New("invalid status string")
}

func (item *ItemStatus) Scan(value interface{}) error {
	bytes, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("fail to scan data from sql: %s", value)

	}

	v, err := parseStr2ItemStatus(string(bytes))

	if err != nil {
		return fmt.Errorf("fail to scan data from sql: %s", value)

	}

	*item = v

	return nil
}

// Convert data type back to JSON
func (item *ItemStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", item.String())), nil
}

// End use enum.

type BookItem struct {
	Id          int         `json:"id" gorm:"column:id;"`
	Name        string      `json:"name" gorm:"column:name;"`
	Description string      `json:"description" gorm:"column:description;"`
	Status      *ItemStatus `json:"status" gorm:"column:status;"`
	CreatedAt   *time.Time  `json:"created_at" gorm:"column:created_at;"`
	UpdatedAt   *time.Time  `json:"updatd_at,omitempty" gorm:"column:updated_at;"`
}

func (BookItem) TableName() string { return "book_items" }

type BookItemCreation struct {
	Id          int    `json:"-" gorm:"column:id;"`
	Name        string `json:"name" gorm:"column:name;"`
	Description string `json:"description" gorm:"column:description;"`
	//Status      string `json:"status" gorm:"column:status;"`
}

func (BookItemCreation) TableName() string { return BookItem{}.TableName() }

type BookItemUpdate struct {
	Name        *string `json:"name" gorm:"column:name;"`
	Description *string `json:"description" gorm:"column:description;"`
	Status      *string `json:"status" gorm:"column:status;"`
}

func (BookItemUpdate) TableName() string { return BookItem{}.TableName() }

type Paging struct {
	Page  int   `json:"page" form:"page"`
	Limit int   `json:"limti" form:"limit"`
	Total int64 `json:"total" form:"-"`
}

func (p *Paging) Process() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
}

func main() {
	// Using GORM for connect db
	dsn := "root:Godpanda01@@tcp(127.0.0.1:3307)/book_list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	// Connect GIN
	r := gin.Default()

	r.Use(cors.Default())

	// CRUD: Create, Read, Update, Delete
	// POST /v1/items (create a new item)
	// GET /v1/items (list items) /v1/items?page=1
	// GET /v1/items/:id (get item detail by id)
	// (PUT || PATCH) /v1/items/:id (update an item by id)
	// DELETE /v1/items/:id (delete item by id)

	//Routes
	r.POST("v1/items", CreateItem(db))
	r.GET("v1/items", ListItem(db))
	r.GET("v1/items/:id", GetItem(db))
	r.PUT("v1/items/:id", UpdateItem(db))
	r.DELETE("v1/items/:id", DeleteItem(db))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8000") //run at port 8000
}

// Create Item
func CreateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data BookItemCreation
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Create(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": data.Id,
		})
	}
}

// Get Item
func GetItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data BookItem

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		//data.Id = id
		if err := db.Where("id = ?", id).First(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}

// Update Item
func UpdateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data BookItemUpdate

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Where("id = ?", id).Updates(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}

// Delete Item
func DeleteItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Table(BookItem{}.TableName()).Where("id = ?", id).Delete(nil).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}

// List Item
func ListItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging Paging
		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		paging.Process()

		var result []BookItem

		//db = db.Where("status <> ?", "In stock")

		if err := db.Table(BookItem{}.TableName()).Count(&paging.Total).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Order("id desc").
			Offset((paging.Page - 1) * paging.Limit).
			Limit(paging.Limit).
			Find(&result).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":   result,
			"paging": paging,
		})
	}
}
