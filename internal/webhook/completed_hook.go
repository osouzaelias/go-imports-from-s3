package webhook

import (
	"bytes"
	"encoding/json"
	"go-import-from-s3/internal"
	"log"
	"net/http"
	"time"
)

const (
	maxRetries   = 3
	initialDelay = time.Second
)

type CompletedHook struct {
	cfg *internal.Config
}

type Payload struct {
	Table          string `json:"table"`
	HashKey        string `json:"hashKey"`
	RangeKey       string `json:"rangeKey"`
	CompletionDate string `json:"completionDate"`
}

func NewCompletedHook(c *internal.Config) *CompletedHook {
	return &CompletedHook{
		cfg: c,
	}
}

func (h CompletedHook) NotifyImportCompleted() error {
	payload := Payload{
		Table:          h.cfg.Table(),
		HashKey:        h.cfg.HashKey(),
		RangeKey:       h.cfg.RangeKey(),
		CompletionDate: time.Now().Format(time.DateTime),
	}

	payloadJSON, _ := json.Marshal(payload)

	var hookError error = nil

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest("POST", h.cfg.Webhook(), bytes.NewBuffer(payloadJSON))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			hookError = err

			log.Println("Erro ao enviar a requisição:", err)

			waitTime := initialDelay * time.Duration(1<<uint(i))

			log.Println("Aguardando a próxima tentativa em:", err)

			time.Sleep(waitTime)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			log.Println("Requisição bem-sucedida!")
			return nil
		} else {
			log.Printf("A requisição retornou status %d, o esperado é 200.", resp.StatusCode)
		}
	}

	return hookError
}
