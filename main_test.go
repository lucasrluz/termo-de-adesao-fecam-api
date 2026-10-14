package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestPostUsuario(t *testing.T) {
	router := setupRouter()

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
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}
