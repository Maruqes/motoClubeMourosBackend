package member

// Aqui defines a estrutura e funções auxiliares
type Member struct {
	id                   int
	numeroSocio          int
	tipoMembro           string
	dataNascimento       string
	dataAdesao           string
	membroResponsavel    string
	nome                 string
	email                string
	telefone             string
	tipoSangue           string
	rua                  string
	numeroPorta          int
	concelho             string
	distrito             string
	codPostal            string
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
