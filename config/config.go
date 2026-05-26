// Guigas pesquisei um pouco e go usa muito de convenções
// privado -> camel case
// publico -> pascal case
// para salvar erros e err
// snak case fortemente não recomendado, me mogaram muito nessa 

packgem config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Blocked []string          // dominios bloqueados
	Words map[string]string  // palavra -> substituto
}

type blockedFile struct {
	Blockeds []string `json:"blockeds"`
}

func LoadFiles(BlockedPath, wordsPath string) *Config {

	conf := &Config{}

	blockedDomainsJSON, err := os.ReadFile(BlockedPath)

	// executa se nao houver erro
	if err == nill {
		var blockedDomains blockedFile

		json.Unmarshal(
            blockedDomainsJSON,
            &blockedDomains,
        )

		config.BlockedDomains = blockedDomains.Domains
	}

	return config
}

// chama o pnteiro em config
// verifica se o domínio recebido existe dentro da lista de domínios bloqueados
func (c *Config) IsBlocked(domain string) bool {
	for _, b := range c.Blocked {
		if string.EqualFold(b, domain) {
			return true
		}
	}
	return false
}