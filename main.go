package main

import (
	"context"
	"fmt"
	"github.com/awakari/int-bluesky/api/grpc/queue"
	"github.com/awakari/int-bluesky/api/http/bluesky"
	"github.com/awakari/int-bluesky/api/http/handler"
	"github.com/awakari/int-bluesky/api/http/interests"
	"github.com/awakari/int-bluesky/api/http/reader"
	"github.com/awakari/int-bluesky/config"
	"github.com/awakari/int-bluesky/model"
	"github.com/awakari/int-bluesky/service/converter"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net/http"
	"os"
)

func main() {

	// init config and logger
	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		panic(fmt.Sprintf("failed to load the config from env: %s", err))
	}
	//
	opts := slog.HandlerOptions{
		Level: slog.Level(cfg.Log.Level),
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, &opts))
	log.Info("starting...")

	svcInterests := interests.NewService(http.DefaultClient, cfg.Api.Interests.Uri, cfg.Api.Token.Internal)
	svcInterests = interests.NewLogging(svcInterests, log)
	log.Info("initialized the Awakari interests API client")

	clientHttp := &http.Client{}

	// init websub reader
	svcReader := reader.NewService(clientHttp, cfg.Api.Reader.Uri)
	svcReader = reader.NewServiceLogging(svcReader, log)
	callbackUrl := fmt.Sprintf(
		"%s://%s:%d%s",
		cfg.Api.Reader.CallBack.Protocol,
		cfg.Api.Reader.CallBack.Host,
		cfg.Api.Reader.CallBack.Port,
		cfg.Api.Reader.CallBack.Path,
	)

	// init queues
	connQueue, err := grpc.NewClient(cfg.Api.Queue.Uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	log.Info("connected to the queue service")
	clientQueue := queue.NewServiceClient(connQueue)
	svcQueue := queue.NewService(clientQueue)
	svcQueue = queue.NewLoggingMiddleware(svcQueue, log)

	err = svcQueue.SetConsumer(context.TODO(), cfg.Api.Queue.InterestsCreated.Name, cfg.Api.Queue.InterestsCreated.Subj)
	if err != nil {
		panic(err)
	}
	log.Info(fmt.Sprintf("initialized the %s queue", cfg.Api.Queue.InterestsCreated.Name))
	go func() {
		err = consumeQueue(
			context.Background(),
			svcReader,
			svcQueue,
			cfg.Api.Queue.InterestsCreated.Name,
			cfg.Api.Queue.InterestsCreated.Subj,
			cfg.Api.Queue.InterestsCreated.BatchSize,
			func(ctx context.Context, svcReader reader.Service, evts []*pb.CloudEvent) {
				consumeInterestEvents(ctx, evts, cfg, log, svcReader, callbackUrl)
			},
		)
		if err != nil {
			panic(err)
		}
	}()

	err = svcQueue.SetConsumer(context.TODO(), cfg.Api.Queue.InterestsUpdated.Name, cfg.Api.Queue.InterestsUpdated.Subj)
	if err != nil {
		panic(err)
	}
	log.Info(fmt.Sprintf("initialized the %s queue", cfg.Api.Queue.InterestsUpdated.Name))
	go func() {
		err = consumeQueue(
			context.Background(),
			svcReader,
			svcQueue,
			cfg.Api.Queue.InterestsUpdated.Name,
			cfg.Api.Queue.InterestsUpdated.Subj,
			cfg.Api.Queue.InterestsUpdated.BatchSize,
			func(ctx context.Context, svcReader reader.Service, evts []*pb.CloudEvent) {
				consumeInterestEvents(ctx, evts, cfg, log, svcReader, callbackUrl)
			},
		)
		if err != nil {
			panic(err)
		}
	}()

	r := gin.Default()

	r.GET("/.well-known/did.json", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
		ctx.JSONP(http.StatusOK, gin.H{
			"@context": "https://www.w3.org/ns/did/v1",
			"id":       fmt.Sprintf("did:web:%s", cfg.Api.Http.Host),
		})
	})

	log.Info(fmt.Sprintf("starting to listen the HTTP API @ port #%d...", cfg.Api.Http.Port))
	go func() {
		err = r.Run(fmt.Sprintf(":%d", cfg.Api.Http.Port))
		if err != nil {
			panic(err)
		}
	}()

	svcConv := converter.NewService(200, true)

	svcBluesky := bluesky.NewService(clientHttp)
	svcBluesky = bluesky.NewLogging(svcBluesky, log)
	did, token, err := svcBluesky.Login(context.Background(), cfg.Api.Bluesky.App.Id, cfg.Api.Bluesky.App.Password)
	if err != nil {
		panic(err)
	}

	hc := handler.NewCallbackHandler(cfg.Api.Reader.Uri, cfg.Api.Http.Host, cfg.Api.EventType, svcConv, svcBluesky, did, token)

	log.Info(fmt.Sprintf("starting to listen the HTTP API @ port #%d...", cfg.Api.Reader.CallBack.Port))
	internalCallbacks := gin.Default()
	internalCallbacks.
		GET(cfg.Api.Reader.CallBack.Path, hc.Confirm).
		POST(cfg.Api.Reader.CallBack.Path, hc.Deliver)
	err = internalCallbacks.Run(fmt.Sprintf(":%d", cfg.Api.Reader.CallBack.Port))
	if err != nil {
		panic(err)
	}
}

func consumeQueue(
	ctx context.Context,
	svcReader reader.Service,
	svcQueue queue.Service,
	name, subj string,
	batchSize uint32,
	consumeEvents func(ctx context.Context, svcReader reader.Service, evts []*pb.CloudEvent),
) (err error) {
	for {
		err = svcQueue.ReceiveMessages(ctx, name, subj, batchSize, func(evts []*pb.CloudEvent) (err error) {
			consumeEvents(ctx, svcReader, evts)
			return
		})
		if err != nil {
			panic(err)
		}
	}
}

func consumeInterestEvents(
	ctx context.Context,
	evts []*pb.CloudEvent,
	cfg config.Config,
	log *slog.Logger,
	svcReader reader.Service,
	callbackUrl string,
) {
	log.Debug(fmt.Sprintf("consumeInterestEvents(%d))\n", len(evts)))
	for _, evt := range evts {

		interestId := evt.GetTextData()
		var groupId string
		if groupIdAttr, groupIdIdPresent := evt.Attributes[model.CeKeyGroupId]; groupIdIdPresent {
			groupId = groupIdAttr.GetCeString()
		}
		if groupId == "" {
			log.Error(fmt.Sprintf("interest %s event: empty group id, skipping", interestId))
			continue
		}

		publicAttr, publicAttrPresent := evt.Attributes[model.CeKeyPublic]
		switch publicAttrPresent && publicAttr.GetCeBoolean() {
		case true:
			_ = svcReader.CreateCallback(ctx, interestId, callbackUrl)
		default:
			log.Debug(fmt.Sprintf("interest %s event: public: %t/%t", interestId, publicAttrPresent, publicAttr.GetCeBoolean()))
		}
	}
	return
}
