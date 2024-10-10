package pg

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"maps"
	"strings"

	"github.com/gowool/cms/internal"
	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

var _ repository.Configuration = (*ConfigurationRepository)(nil)

const (
	cfgSelectSQL = "SELECT key,value FROM pages_configuration"
	cfgInsertSQL = "INSERT INTO pages_configuration (key,value) VALUES %s ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value"
)

type ConfigurationRepository struct {
	db *sql.DB
}

func NewConfigurationRepository(db *sql.DB) *ConfigurationRepository {
	return &ConfigurationRepository{db: db}
}

func (r *ConfigurationRepository) Load(ctx context.Context) (model.Configuration, error) {
	rows, err := r.db.QueryContext(ctx, cfgSelectSQL)
	if err != nil {
		return model.Configuration{}, err
	}
	defer func() {
		_ = rows.Close()
	}()

	m := model.NewConfiguration()

	for rows.Next() {
		var key, value string
		if err = rows.Scan(&key, &value); err != nil {
			return model.Configuration{}, err
		}

		switch key {
		case "debug":
			m.Debug = value == "true"
		case "multisite":
			m.Multisite = model.MultisiteStrategy(value)
		case "ignore_request_patterns":
			if err = json.Unmarshal(internal.Bytes(value), &m.IgnoreRequestPatterns); err != nil {
				return model.Configuration{}, err
			}
		case "ignore_request_uris":
			if err = json.Unmarshal(internal.Bytes(value), &m.IgnoreRequestURIs); err != nil {
				return model.Configuration{}, err
			}
		case "fallback_locale":
			m.FallbackLocale = value
		case "catch_errors":
			if err = json.Unmarshal(internal.Bytes(value), &m.CatchErrors); err != nil {
				return model.Configuration{}, err
			}
		default:
			m.Additional[key] = value
		}
	}

	return m, nil
}

func (r *ConfigurationRepository) Save(ctx context.Context, m *model.Configuration) error {
	data := toMap(m)

	values := make([]string, 0, len(data))
	args := make([]any, 0, len(data)*2)

	for k, v := range data {
		args = append(args, k, v)
		values = append(values, fmt.Sprintf("($%d,$%d)", len(args)-1, len(args)))
	}

	query := fmt.Sprintf(cfgInsertSQL, strings.Join(values, ","))

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func toMap(cfg *model.Configuration) (data map[string]string) {
	if cfg.Additional == nil {
		data = make(map[string]string)
	} else {
		data = maps.Clone(cfg.Additional)
	}

	data["debug"] = fmt.Sprintf("%t", cfg.Debug)
	data["multisite"] = cfg.Multisite.String()
	data["ignore_request_patterns"] = jsonString(cfg.IgnoreRequestPatterns)
	data["ignore_request_uris"] = jsonString(cfg.IgnoreRequestURIs)
	data["fallback_locale"] = cfg.FallbackLocale
	data["catch_errors"] = jsonString(cfg.CatchErrors)
	return data
}

func jsonString(data any) string {
	raw, _ := json.Marshal(data)
	return internal.String(raw)
}
