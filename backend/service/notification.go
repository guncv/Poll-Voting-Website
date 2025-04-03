package service

import (
	"context"

	"github.com/guncv/Poll-Voting-Website/backend/entity"
	"github.com/guncv/Poll-Voting-Website/backend/log"
	"github.com/guncv/Poll-Voting-Website/backend/repository"
)

type INotificationService interface {
	SendAlertReachScoreToAdmin(ctx context.Context, subject, message string) error
	NotifyUserOfUrgentQuestion(ctx context.Context, subject, message string) error
	AddSubscriberToUserTopic(ctx context.Context, email string) error
}

type NotificationService struct {
	sender repository.INotificationRepository
	log    log.LoggerInterface
}

func NewNotificationService(sender repository.INotificationRepository, logger log.LoggerInterface) INotificationService {
	return &NotificationService{
		sender: sender,
		log:    logger,
	}
}

func (a *NotificationService) SendAlertReachScoreToAdmin(ctx context.Context, subject, message string) error {
	a.log.InfoWithID(ctx, "[Service: SendAlertReachScoreToAdmin] Called")
	alert := entity.Alert{
		Subject: subject,
		Message: message,
	}

	err := a.sender.SendAdminAlert(ctx, alert)
	if err != nil {
		a.log.ErrorWithID(ctx, "[Service: SendAlertReachScoreToAdmin] Failed to send alert:", err)
		return err
	}
	a.log.InfoWithID(ctx, "[Service: SendAlertReachScoreToAdmin] Alert sent to SNS!")
	return nil
}

func (a *NotificationService) NotifyUserOfUrgentQuestion(ctx context.Context, subject, message string) error {
	a.log.InfoWithID(ctx, "[Service: NotifyUserOfUrgentQuestion] Called")
	alert := entity.Alert{
		Subject: subject,
		Message: message,
	}
	err := a.sender.SendUserAlert(ctx, alert)
	if err != nil {
		a.log.ErrorWithID(ctx, "[Service: NotifyUserOfUrgentQuestion] Failed to send alert:", err)
		return err
	}
	return nil
}

func (a *NotificationService) AddSubscriberToUserTopic(ctx context.Context, email string) error {
	a.log.InfoWithID(ctx, "[Service: AddSubscriberToUserTopic] Called")
	err := a.sender.SubscribeToUserTopic(ctx, email)
	if err != nil {
		a.log.ErrorWithID(ctx, "[Service: AddSubscriberToUserTopic] Failed to subscribe:", err)
		return err
	}
	return nil
}
