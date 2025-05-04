package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/awakari/int-bluesky/api/http/bluesky"
	"github.com/awakari/int-bluesky/api/http/pub"
	"github.com/awakari/int-bluesky/api/http/reader"
	"github.com/awakari/int-bluesky/config"
	"github.com/awakari/int-bluesky/model"
	"github.com/awakari/int-bluesky/service/converter"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
)

type Service interface {
	ConsumeInterestEvents(ctx context.Context, evts []*pb.CloudEvent) (err error)
	ConsumePostEvents(ctx context.Context, evts []*pb.CloudEvent) (err error)
}

type service struct {
	cfg         config.Config
	svcReader   reader.Service
	callbackUrl string
	svcConv     converter.Service
	svcPub      pub.Service
	svcBluesky  bluesky.Service
	didWeb      string
	didPlc      string
	token       string
}

func NewService(
	cfg config.Config,
	svcReader reader.Service,
	callbackUrl string,
	svcConv converter.Service,
	svcPub pub.Service,
	svcBluesky bluesky.Service,
	didWeb string,
	didPlc string,
	token string,
) Service {
	return service{
		cfg:         cfg,
		svcReader:   svcReader,
		callbackUrl: callbackUrl,
		svcConv:     svcConv,
		svcPub:      svcPub,
		svcBluesky:  svcBluesky,
		didWeb:      didWeb,
		didPlc:      didPlc,
		token:       token,
	}
}

func (s service) ConsumeInterestEvents(ctx context.Context, evts []*pb.CloudEvent) (err error) {
	for _, evt := range evts {
		interestId := evt.GetTextData()
		var groupId string
		if groupIdAttr, groupIdIdPresent := evt.Attributes[model.CeKeyGroupId]; groupIdIdPresent {
			groupId = groupIdAttr.GetCeString()
		}
		if groupId == "" {
			err = errors.Join(err, fmt.Errorf("interest %s event: empty group id, skipping", interestId))
			continue
		}
		publicAttr, publicAttrPresent := evt.Attributes[model.CeKeyPublic]
		if publicAttrPresent && publicAttr.GetCeBoolean() {
			_ = s.svcReader.CreateCallback(ctx, interestId, s.callbackUrl)
			_ = s.svcBluesky.CreateFeed(ctx, s.didWeb, s.didPlc, s.token, interestId)
		}
	}
	return
}

func (s service) ConsumePostEvents(ctx context.Context, evts []*pb.CloudEvent) (err error) {
	for _, src := range evts {
		dst, userId, errConv := s.svcConv.PostToEvent(ctx, src)
		if errConv != nil {
			err = errors.Join(err, errConv)
			continue
		}
		if dst == nil {
			continue // may be nil if the source event is not a bluesky post
		}
		if errPub := s.svcPub.Publish(ctx, dst, s.cfg.Api.GroupId, userId); errPub != nil {
			err = errors.Join(err, errPub)
		}
	}
	return
}
