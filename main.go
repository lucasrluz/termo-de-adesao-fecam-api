package main

import (
	"fmt"
	"log"
	"os"
	"termo-de-adesao-fecam-api/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SolicitacaoDTO struct {
	Entidade              string `json:"entidade"`
	Sistemas              []uint `json:"sistemas"`
	RepresentanteNome     string `json:"representante_nome"`
	RepresentanteCPF      string `json:"representante_cpf"`
	RepresentanteCargo    string `json:"representante_cargo"`
	RepresentanteEmail    string `json:"representante_email"`
	RepresentanteTelefone string `json:"representante_telefone"`
	ResponsavelNome       string `json:"responsavel_nome"`
	ResponsavelCPF        string `json:"responsavel_cpf"`
	ResponsavelCargo      string `json:"responsavel_cargo"`
	ResponsavelEmail      string `json:"responsavel_email"`
	ResponsavelTelefone   string `json:"responsavel_telefone"`
}

func setupDB() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&models.Solicitacao{})

	return db
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	db := setupDB()

	router.POST("/api/solicitacao", func(c *gin.Context) {
		var solicitacao SolicitacaoDTO

		if err := c.ShouldBindJSON(&solicitacao); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		solicitacaoModel := models.Solicitacao{
			Status:                models.StatusPendente,
			Entidade:              solicitacao.Entidade,
			RepresentanteNome:     solicitacao.RepresentanteNome,
			RepresentanteCPF:      solicitacao.RepresentanteCPF,
			RepresentanteCargo:    solicitacao.RepresentanteCargo,
			RepresentanteEmail:    solicitacao.RepresentanteEmail,
			RepresentanteTelefone: solicitacao.RepresentanteTelefone,
			ResponsavelNome:       solicitacao.ResponsavelNome,
			ResponsavelCPF:        solicitacao.ResponsavelCPF,
			ResponsavelCargo:      solicitacao.ResponsavelCargo,
			ResponsavelEmail:      solicitacao.ResponsavelEmail,
			ResponsavelTelefone:   solicitacao.ResponsavelTelefone,
		}

		for _, id := range solicitacao.Sistemas {
			solicitacaoModel.Servicos = append(solicitacaoModel.Servicos, models.Servico{
				IDServico: id,
			})
		}

		db.Create(&solicitacaoModel)

		c.JSON(201, solicitacao)
	})

	return router
}

func main() {
	router := setupRouter()
	router.Run(":8080")
}
