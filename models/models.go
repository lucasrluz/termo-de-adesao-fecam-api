package models

import "time"

// StatusSolicitacao representa os estados possíveis de uma solicitação.
type StatusSolicitacao string

const (
	StatusPendente StatusSolicitacao = "pendente"
	StatusNegado   StatusSolicitacao = "negado"
	StatusAprovado StatusSolicitacao = "aprovado"
)

// Solicitacao -> tb_solicitacao
type Solicitacao struct {
	IDSolicitacao uint              `gorm:"column:id_solicitacao;primaryKey;autoIncrement"`
	Status        StatusSolicitacao `gorm:"column:status;type:varchar(20);not null;default:'pendente';index"`
	Entidade      string            `gorm:"column:entidade;type:varchar(255);not null"`

	RepresentanteNome     string `gorm:"column:representante_nome;type:varchar(255);not null"`
	RepresentanteCPF      string `gorm:"column:representante_cpf;type:varchar(14);not null;index"`
	RepresentanteEmail    string `gorm:"column:representante_email;type:varchar(255);not null"`
	RepresentanteCargo    string `gorm:"column:representante_cargo;type:varchar(100)"`
	RepresentanteTelefone string `gorm:"column:representante_telefone;type:varchar(20)"`

	ResponsavelNome     string `gorm:"column:responsavel_nome;type:varchar(255);not null"`
	ResponsavelCPF      string `gorm:"column:responsavel_cpf;type:varchar(14);not null;index"`
	ResponsavelEmail    string `gorm:"column:responsavel_email;type:varchar(255);not null"`
	ResponsavelCargo    string `gorm:"column:responsavel_cargo;type:varchar(100)"`
	ResponsavelTelefone string `gorm:"column:responsavel_telefone;type:varchar(20)"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// 1:N -> uma solicitação possui vários registros na tabela de junção
	Servicos []Servico `gorm:"many2many:tb_solicitacao_servicos;"`
}

func (Solicitacao) TableName() string { return "tb_solicitacao" }

// Servico -> tb_servico
type Servico struct {
	IDServico uint   `gorm:"column:id_servico;primaryKey;autoIncrement"`
	Nome      string `gorm:"column:nome;type:varchar(255);not null;uniqueIndex"`
}

func (Servico) TableName() string { return "tb_servico" }
