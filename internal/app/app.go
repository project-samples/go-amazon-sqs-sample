package app

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/core-go/health"
	hm "github.com/core-go/health/mongo"
	w "github.com/core-go/mongo/writer"
	"github.com/core-go/mq"
	v "github.com/core-go/mq/validator"
	"github.com/core-go/mq/zap"
	"github.com/core-go/sqs"
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	Receive       func(ctx context.Context, handle func(context.Context, []byte, map[string]string))
	Handle        func(context.Context, []byte, map[string]string)
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	log.Initialize(cfg.Log)
	client, er0 := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.Uri))
	if er0 != nil {
		log.Error(ctx, "Cannot connect to MongoDB: Error: "+er0.Error())
		return nil, er0
	}
	db := client.Database(cfg.Mongo.Database)

	logError := log.ErrorMsg
	var logInfo func(context.Context, string)
	if log.IsInfoEnable() {
		logInfo = log.InfoMsg
	}

	sqsClient, er1 := sqs.Connect(cfg.SQS)
	if er1 != nil {
		log.Error(ctx, "Cannot create a new sqs. Error: "+er1.Error())
		return nil, er1
	}

	receiver, er2 := sqs.NewReceiverByQueueName(sqsClient, cfg.SQS.QueueName, true, 20, 1)
	if er2 != nil {
		log.Error(ctx, "Cannot create a new receiver. Error: "+er2.Error())
		return nil, er2
	}
	validator, er3 := v.NewValidator[*User]()
	if er3 != nil {
		return nil, er3
	}
	errorHandler := mq.NewErrorHandler[*User](logError)
	sender, er4 := sqs.NewSenderByQueueName(sqsClient, cfg.SQS.QueueName, 1)
	if er4 != nil {
		return nil, er4
	}
	writer := w.NewWriter[*User](db, "user")
	han := mq.NewRetryHandlerByConfig[User](cfg.Retry, writer.Write, validator.Validate, errorHandler.RejectWithMap, nil, sender.Send, logError, logInfo)
	mongoChecker := hm.NewHealthChecker(client)
	receiverChecker := sqs.NewHealthChecker(sqsClient, cfg.SQS.QueueName)
	healthHandler := health.NewHandler(mongoChecker, receiverChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		Receive:       receiver.Receive,
		Handle:        han.Handle,
	}, nil
}
