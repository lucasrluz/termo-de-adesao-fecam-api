package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"gorm.io/gorm"
)

var (
	testDB     *gorm.DB
	testRouter *gin.Engine
)

func truncateAll(db *gorm.DB) error {
	return db.Exec(`TRUNCATE TABLE tb_solicitacao, tb_servico, tb_solicitacao_servicos RESTART IDENTITY CASCADE`).Error
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	testDB = setupDB()

	if err := truncateAll(testDB); err != nil {
		log.Fatalf("falha ao limpar o banco: %v", err)
	}

	testRouter = setupRouter(testDB)

	os.Exit(m.Run())
}

func TestPostUsuario(t *testing.T) {
	w := httptest.NewRecorder()

	exampleUser := SolicitacaoDTO{
		Entidade:              "Exemplo Entidade",
		Sistemas:              []uint{1},
		RepresentanteNome:     "João da Silva",
		RepresentanteCPF:      "123.456.789-00",
		RepresentanteCargo:    "Gerente",
		RepresentanteEmail:    "joao.silva@example.com",
		RepresentanteTelefone: "(11) 98765-4321",
		ResponsavelNome:       "Maria Oliveira",
		ResponsavelCPF:        "987.654.321-00",
		ResponsavelCargo:      "Coordenadora",
		ResponsavelEmail:      "maria.oliveira@example.com",
		ResponsavelTelefone:   "(11) 98765-4321",
	}

	userJson, _ := json.Marshal(exampleUser)
	req, _ := http.NewRequest("POST", "/api/solicitacao", strings.NewReader(string(userJson)))
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}
