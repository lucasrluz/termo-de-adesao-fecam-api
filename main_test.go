package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"termo-de-adesao-fecam-api/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	testDB     *gorm.DB
	testRouter *gin.Engine
)

func truncateAll(db *gorm.DB) error {
	return db.Exec(`TRUNCATE TABLE tb_solicitacao, tb_solicitacao_servicos RESTART IDENTITY CASCADE`).Error
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
	require.NoError(t, truncateAll(testDB))

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

	// consulta no banco: deve existir exatamente uma solicitação
	var total int64
	require.NoError(t, testDB.Model(&models.Solicitacao{}).Count(&total).Error)
	require.EqualValues(t, 1, total)

	var salva models.Solicitacao
	require.NoError(t, testDB.Preload("Servicos").First(&salva).Error)

	assert.Equal(t, models.StatusPendente, salva.Status) // default do model
	assert.Equal(t, exampleUser.Entidade, salva.Entidade)

	assert.Equal(t, exampleUser.RepresentanteNome, salva.RepresentanteNome)
	assert.Equal(t, exampleUser.RepresentanteCPF, salva.RepresentanteCPF)
	assert.Equal(t, exampleUser.RepresentanteCargo, salva.RepresentanteCargo)
	assert.Equal(t, exampleUser.RepresentanteEmail, salva.RepresentanteEmail)
	assert.Equal(t, exampleUser.RepresentanteTelefone, salva.RepresentanteTelefone)

	assert.Equal(t, exampleUser.ResponsavelNome, salva.ResponsavelNome)
	assert.Equal(t, exampleUser.ResponsavelCPF, salva.ResponsavelCPF)
	assert.Equal(t, exampleUser.ResponsavelCargo, salva.ResponsavelCargo)
	assert.Equal(t, exampleUser.ResponsavelEmail, salva.ResponsavelEmail)
	assert.Equal(t, exampleUser.ResponsavelTelefone, salva.ResponsavelTelefone)

	// associação many2many
	require.Len(t, salva.Servicos, 1)
	assert.Equal(t, uint(1), salva.Servicos[0].IDServico)
	assert.Equal(t, "Portal Municipal", salva.Servicos[0].Nome)
}
