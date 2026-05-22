package server

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Importar os docs gerados pelo Swagger
	_ "simple-erp-service/docs"
)

// SetupSwagger configura o Swagger na aplicação
func (s *Server) setupSwagger() {
	// URL para acessar a documentação Swagger
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
