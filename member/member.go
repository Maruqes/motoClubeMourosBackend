package member

// Aqui defines a estrutura e funções auxiliares
type Member struct {
	ID                   string
	NumeroSocio          int
	Junior               bool
	SocioResponsavel     string
	DataNascimento       string
	DataAdesao           string
	MembroResponsavel    string
	Nome                 string
	Email                string
	Telefone             string
	TipoSangue           string
	Rua                  string
	Numero               string
	Concelho             string
	Distrito             string
	CodPostal            string
	Tipo                 string
	GrupoWA              bool
	DataInscricao        string
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func boolPtrToInt(b *bool) int {
	if b == nil {
		return 0
	}
	if *b {
		return 1
	}
	return 0
}
