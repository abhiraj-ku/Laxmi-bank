package api

import (
	db "example/laxmi_chit_fund/db/sqlc"
	"example/laxmi_chit_fund/utils"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config utils.Config
	store  db.Store
	router *gin.Engine
}

// func NewServer(store *db.Store) *Server{

// }

// Setup Gin Router to handle all the incoming API routes

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createAccount)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}

}
