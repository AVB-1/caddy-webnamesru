package webnamesru

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/libdns/libdns"
)

const apiBase = "https://www.webnames.ru/scripts/json_domain_zone_manager.pl"

type Provider struct {
	APIKey string `json:"api_key,omitempty"`
}

// AppendRecords добавляет TXT-запись для DNS-01 challenge.
func (p *Provider) AppendRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	var added []libdns.Record
	for _, rec := range recs {
		txt, ok := rec.(libdns.TXT)
		if !ok {
			return nil, fmt.Errorf("webnamesru: поддерживаются только TXT-записи")
		}

		// libdns передаёт относительное имя (например, "_acme-challenge").
		// webnames.ru API ожидает "subDomain:value" в параметре record.
		recordParam := fmt.Sprintf("%s:%s", txt.Name, txt.Text)
		if err := p.doRequest(ctx, zone, "add", recordParam); err != nil {
			return added, err
		}
		added = append(added, rec)
	}
	return added, nil
}

// DeleteRecords удаляет TXT-запись после проверки.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, recs []libdns.Record) error {
	for _, rec := range recs {
		txt, ok := rec.(libdns.TXT)
		if !ok {
			continue
		}
		recordParam := fmt.Sprintf("%s:%s", txt.Name, txt.Text)
		if err := p.doRequest(ctx, zone, "delete", recordParam); err != nil {
			return err
		}
		// Небольшая пауза, чтобы DNS успел обновиться (как делает lego).
		time.Sleep(2 * time.Second)
	}
	return nil
}

func (p *Provider) doRequest(ctx context.Context, domain, action, record string) error {
	u, err := url.Parse(apiBase)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("apikey", p.APIKey)
	q.Set("domain", domain)
	q.Set("type", "TXT")
	q.Set("record", record)
	q.Set("action", action)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("webnamesru: не удалось создать запрос: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("webnamesru: HTTP запрос не удался: %w", err)
	}
	defer resp.Body.Close()

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return fmt.Errorf("webnamesru: не удалось разобрать ответ: %w", err)
	}

	if result, ok := body["result"].(string); !ok || result != "OK" {
		return fmt.Errorf("webnamesru: API вернул ошибку (action=%s, record=%s): %v", action, record, body["result"])
	}
	return nil
}
