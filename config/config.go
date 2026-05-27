// Guigas pesquisei um pouco e go usa muito de convenções
// privado -> camel case
// publico -> pascal case
// para salvar erros e err
// snak case fortemente não recomendado, me mogaram muito nessa
package config

import (
	"encoding/json"
	"os"
	"strings"
)

type Config struct {
	Blocked []string          // dominios bloqueados
	Words   map[string]string // palavra -> substituto
}

type blockedFile struct {
	Bloqueados []string `json:"bloqueados"`
}

func LoadFiles(blockedPath, wordsPath string) *Config {
	config := &Config{}

	blockedDomainsJSON, err := os.ReadFile(blockedPath)
	// executa se nao houver erro
	if err == nil {
		var blockedDomains blockedFile
		json.Unmarshal(
			blockedDomainsJSON,
			&blockedDomains,
		)
		config.Blocked = blockedDomains.Bloqueados
	}

	// carrega words.json
	wordsJSON, err := os.ReadFile(wordsPath)
	if err == nil {
		json.Unmarshal(wordsJSON, &config.Words)
	}

	return config
}

// chama o ponteiro em config
// verifica se o domínio recebido existe dentro da lista de domínios bloqueados
func (c *Config) IsBlocked(domain string) bool {
	for _, b := range c.Blocked {
		if strings.EqualFold(b, domain) {
			return true
		}
	}
	return false
}
