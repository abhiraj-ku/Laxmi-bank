package api

import (
	db "example/laxmi_chit_fund/db/sqlc"

	"github.com/gin-gonic/gin"
)

// servess
type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: *store}
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	server.router = router
	return server
}

// Setup Gin Router to handle all the incoming API routes

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createAccount)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}

}
