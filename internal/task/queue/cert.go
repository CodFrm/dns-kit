package queue

import (
	"context"

	"github.com/codfrm/cago/pkg/broker"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/dns-kit/internal/task/queue/message"
)

const (
	CertCreate = "cert.create"
)

func PublishCertCreate(ctx context.Context, message *message.CreateCertMessage) error {
	return broker.Default().Publish(ctx, CertCreate, &broker2.Message{
		Body: message.Marshal(),
	})
}

func SubscribeCertCreate(ctx context.Context,
	handler func(ctx context.Context, msg *message.CreateCertMessage) error, opts ...broker2.SubscribeOption) error {
	_, err := broker.Default().Subscribe(ctx, CertCreate, func(ctx context.Context, msg broker2.Event) error {
		m := &message.CreateCertMessage{}
		if err := m.Unmarshal(msg.Message().Body); err != nil {
			return err
		}
		return handler(ctx, m)
	}, opts...)
	return err
}
