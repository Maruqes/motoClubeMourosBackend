package events


// Aqui defines a estrutura e funções auxiliares
type Event struct {
	id                   int
	nome	             string
	descricao            string
	dataEvento		     string
	dataLimiteInscricao  string
	numInscricoes    	 int
	participantes		 string // 2,4,65,23,12
	cartaz 				 string // @/images/cartaz.jpeg	
	programa			 string
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
