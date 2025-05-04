package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	Api ApiConfig
	Log struct {
		Level int `envconfig:"LOG_LEVEL" default:"-4" required:"true"`
	}
}

type ApiConfig struct {
	Bluesky struct {
		App struct {
			Id       string `required:"true" envconfig:"API_BLUESKY_APP_ID" default:"bluesky.awakari.com"`
			Password string `required:"true" envconfig:"API_BLUESKY_APP_PASSWORD"`
		}
	}
	Http struct {
		Host string `envconfig:"API_HTTP_HOST" required:"true"`
		Port uint16 `envconfig:"API_HTTP_PORT" default:"8080" required:"true"`
	}
	Metrics struct {
		Port uint16 `envconfig:"API_METRICS_PORT" default:"9090" required:"true"`
	}
	EventType EventTypeConfig
	Interests struct {
		Uri              string `envconfig:"API_INTERESTS_URI" required:"true" default:"http://interests-api:8080/v1"`
		DetailsUriPrefix string `envconfig:"API_INTERESTS_DETAILS_URI_PREFIX" required:"true" default:"https://awakari.com/sub-details.html?id="`
		Grpc             struct {
			Uri        string `envconfig:"API_INTERESTS_GRPC_URI" default:"interests-api:50051" required:"true"`
			Connection struct {
				Count struct {
					Init uint32 `envconfig:"API_INTERESTS_GRPC_CONN_COUNT_INIT" default:"1" required:"true"`
					Max  uint32 `envconfig:"API_INTERESTS_GRPC_CONN_COUNT_MAX" default:"5" required:"true"`
				}
				IdleTimeout time.Duration `envconfig:"API_INTERESTS_GRCP_CONN_IDLE_TIMEOUT" default:"15m" required:"true"`
			}
		}
	}
	Writer struct {
		Backoff time.Duration `envconfig:"API_WRITER_BACKOFF" default:"10s" required:"true"`
		Timeout time.Duration `envconfig:"API_WRITER_TIMEOUT" default:"10s" required:"true"`
		Uri     string        `envconfig:"API_WRITER_URI" default:"http://pub:8080/v1" required:"true"`
	}
	Reader     ReaderConfig
	Prometheus PrometheusConfig
	Queue      QueueConfig
	Token      struct {
		Internal string `envconfig:"API_TOKEN_INTERNAL" required:"true"`
	}
	GroupId string `envconfig:"API_GROUP_ID" default:"default" required:"true"`
}

type PrometheusConfig struct {
	Uri string `envconfig:"API_PROMETHEUS_URI" default:"http://prometheus-server:80" required:"true"`
}

type ReaderConfig struct {
	Uri          string `envconfig:"API_READER_URI" default:"http://reader:8080/v1" required:"true"`
	UriEventBase string `envconfig:"API_READER_URI_EVT_BASE" default:"https://awakari.com/pub-msg.html?id=" required:"true"`
	CallBack     struct {
		Protocol string `envconfig:"API_READER_CALLBACK_PROTOCOL" default:"http" required:"true"`
		Host     string `envconfig:"API_READER_CALLBACK_HOST" default:"int-bluesky" required:"true"`
		Port     uint16 `envconfig:"API_READER_CALLBACK_PORT" default:"8081" required:"true"`
		Path     string `envconfig:"API_READER_CALLBACK_PATH" default:"/v1/callback" required:"true"`
	}
}

type EventTypeConfig struct {
	InterestsUpdated string `envconfig:"API_EVENT_TYPE_INTERESTS_UPDATED" required:"true" default:"interests-updated"`
}

type QueueConfig struct {
	Uri              string `envconfig:"API_QUEUE_URI" default:"queue:50051" required:"true"`
	InterestsCreated struct {
		BatchSize uint32 `envconfig:"API_QUEUE_INTERESTS_CREATED_BATCH_SIZE" default:"1" required:"true"`
		Name      string `envconfig:"API_QUEUE_INTERESTS_CREATED_NAME" default:"int-bluesky" required:"true"`
		Subj      string `envconfig:"API_QUEUE_INTERESTS_CREATED_SUBJ" default:"interests-created" required:"true"`
	}
	InterestsUpdated struct {
		BatchSize uint32 `envconfig:"API_QUEUE_INTERESTS_UPDATED_BATCH_SIZE" default:"1" required:"true"`
		Name      string `envconfig:"API_QUEUE_INTERESTS_UPDATED_NAME" default:"int-bluesky" required:"true"`
		Subj      string `envconfig:"API_QUEUE_INTERESTS_UPDATED_SUBJ" default:"interests-updated" required:"true"`
	}
	SourceWebsocket struct {
		BatchSize uint32 `envconfig:"API_QUEUE_SRC_WEBSOCKET_BATCH_SIZE" default:"100" required:"true"`
		Name      string `envconfig:"API_QUEUE_SRC_WEBSOCKET_NAME" default:"int-bluesky" required:"true"`
		Subj      string `envconfig:"API_QUEUE_SRC_WEBSOCKET_SUBJ" default:"source-websocket-bluesky" required:"true"`
	}
}

func NewConfigFromEnv() (cfg Config, err error) {
	err = envconfig.Process("", &cfg)
	return
}
