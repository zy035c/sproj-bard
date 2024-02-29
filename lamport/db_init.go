package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func db_init() {
	// init db
	// create table
	// return db

	// db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	// TODO
}

func make_conn(c *gin.Context) *gorm.DB {
	db, _ := c.Get("db_conn")
	db_conn := db.(*gorm.DB)
	return db_conn
}
