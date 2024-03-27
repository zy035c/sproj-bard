package controller

import (
	"fmt"
	"lamport/models"

	"github.com/gin-gonic/gin"
)

func make_fake_db(c *gin.Context) *map[uint64]models.Order {
	// create a fake db
	// return db
	db, _ := c.Get("fake_db")
	db_conn := db.(*map[uint64]models.Order)
	return db_conn
}

func InsertOrUpdate(c *gin.Context, order *models.Order, timestamp uint64) (*models.Order, error) {

	db := make_fake_db(c)

	order.Timestamp = timestamp
	// insert
	(*db)[timestamp] = *order

	fmt.Println("[db] Inserted: ", *order)

	return order, nil

	// db := make_conn(c)
	// ord, err := models.GetMostRecent(timestamp, db)
	// if err != nil {
	// 	c.JSON(404, gin.H{"error": err})
	// 	return nil, err
	// }

	// return ord, nil

}

func GetMostRecent(c *gin.Context, timestamp uint64) (*models.Order, error) {
	// get by given timestamp
	// return the struct
	db := make_fake_db(c)
	if order, ok := (*db)[timestamp]; ok {
		return &order, nil
	}
	return nil, fmt.Errorf("NotFoundError timestamp: %d", timestamp)
}

func GetLastProcID(c *gin.Context, timestamp uint64) int {
	// get the last proc id
	// return the id
	db := make_fake_db(c)
	if order, ok := (*db)[timestamp]; ok {
		return order.ProcID
	}
	return -1
}
