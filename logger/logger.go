// registra todas as requisições recebidas pelo proxy

package logger

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Entry struct {
	Timestamp string `json:"timestamp"`
	URL       string `json:"url"`
	Action    string `json:"action"` // permitido bloqueado filtrado
}

type Logger struct {
	mutex   sync.Mutex // apenas uma goroutine vai escrever o log por vez
	path    string
	entries []Entry
}

// Cria um logger
func New(path string) *Logger {

	l := &Logger{path: path}
	data, err := os.ReadFile(path)

	if err == nil {
		json.Unmarshal(data, &l.entries)
	}

	return l
}

func (l *Logger) Write(rawURL, action string) {
	entry := Entry{
		Timestamp: time.Now().Format(time.RFC3339),
		URL:       rawURL,
		Action:    action,
	}

	l.mutex.Lock()
	// adquire o lock, outras goroutines ficam em fila aqui
	defer l.mutex.Unlock() // libera o lock ao sair da funcao

	l.entries = append(l.entries, entry)

	data, _ := json.MarshalIndent(l.entries, "", " ")
	os.WriteFile(l.path, data, 0644)
}
