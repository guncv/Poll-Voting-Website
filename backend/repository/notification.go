package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	cfg "github.com/guncv/Poll-Voting-Website/backend/config"
	"github.com/guncv/Poll-Voting-Website/backend/entity"
	"github.com/guncv/Poll-Voting-Website/backend/log"
)

type NotificationRepository struct {
	client *sns.Client
	cfg    cfg.Config
	log    log.LoggerInterface
}

type INotificationRepository interface {
	SendAdminAlert(ctx context.Context, alert entity.Alert) error
	SendUserAlert(ctx context.Context, alert entity.Alert) error
	SubscribeToUserTopic(ctx context.Context, email string) error
}

func NewNotificationRepository(client *sns.Client, cfg cfg.Config, logger log.LoggerInterface) INotificationRepository {
	return &NotificationRepository{
		client: client,
		cfg:    cfg,
		log:    logger,
	}
}

func (s *NotificationRepository) SendAdminAlert(ctx context.Context, alert entity.Alert) error {
	s.log.InfoWithID(ctx, "[Repository: SendAdminAlert] Called")
	return s.publishAlert(ctx, s.cfg.Notification.AdminTopicArn, "[Repository: SendAdminAlert]", alert)
}

func (s *NotificationRepository) SendUserAlert(ctx context.Context, alert entity.Alert) error {
	s.log.InfoWithID(ctx, "[Repository: SendUserAlert] Called")
	return s.publishAlert(ctx, s.cfg.Notification.UserTopicArn, "[Repository: SendUserAlert]", alert)
}

func (s *NotificationRepository) SubscribeToUserTopic(ctx context.Context, email string) error {
	s.log.InfoWithID(ctx, "[Repository: SubscribeToUserTopic] Called")
	return s.subscribeEmail(ctx, s.cfg.Notification.UserTopicArn, "[Repository: SubscribeToUserTopic]", email)
}

func (s *NotificationRepository) publishAlert(ctx context.Context, topicArn, logPrefix string, alert entity.Alert) error {
	s.log.InfoWithID(ctx, logPrefix+" Called")

	_, err := s.client.Publish(ctx, &sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Subject:  aws.String(alert.Subject),
		Message:  aws.String(alert.Message),
	})

	if err != nil {
		s.log.ErrorWithID(ctx, logPrefix+" Failed to publish message:", err)
		return fmt.Errorf("failed to publish message to topic %s: %w", topicArn, err)
	}

	s.log.InfoWithID(ctx, logPrefix+" Alert successfully published")
	return nil
}

func (s *NotificationRepository) subscribeEmail(ctx context.Context, topicArn, logPrefix, email string) error {
	s.log.InfoWithID(ctx, logPrefix+" Called")

	_, err := s.client.Subscribe(ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicArn),
		Protocol: aws.String("email"),
		Endpoint: aws.String(email),
	})

	if err != nil {
		s.log.ErrorWithID(ctx, logPrefix+" Failed to subscribe:", err)
		return fmt.Errorf("failed to subscribe %s to topic %s: %w", email, topicArn, err)
	}

	s.log.InfoWithID(ctx, logPrefix+" Subscription requested for "+email)
	return nil
}
