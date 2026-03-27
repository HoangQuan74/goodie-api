package elasticsearch

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/kainguyen/goodie-api/pkg/logger"
	"go.uber.org/zap"
)

type Config struct {
	URL string
}

func NewClient(cfg Config) (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.URL},
	})
	if err != nil {
		return nil, fmt.Errorf("create elasticsearch client: %w", err)
	}

	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("ping elasticsearch: %w", err)
	}
	defer res.Body.Close()

	logger.Get().Info("connected to Elasticsearch",
		zap.String("url", cfg.URL),
	)

	return es, nil
}
