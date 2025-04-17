package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/guncv/Poll-Voting-Website/backend/entity"
	"github.com/guncv/Poll-Voting-Website/backend/log"
	"github.com/guncv/Poll-Voting-Website/backend/repository"
)

type INotificationService interface {
	SendAlertReachParticipantsToAdmin(ctx context.Context, questionText string, totalParticipants int, firstChoice, secondChoice string, firstChoiceCount, secondChoiceCount int) error
	NotifyUserOfAdminQuestion(ctx context.Context, email, subject, message string) error
	AddSubscriberToUserTopic(ctx context.Context, email string) error
	CheckIsAdmin(ctx context.Context, email string) (bool, error)
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

func (a *NotificationService) SendAlertReachParticipantsToAdmin(ctx context.Context, questionText string, totalParticipants int, firstChoice, secondChoice string, firstChoiceCount, secondChoiceCount int) error {
	a.log.InfoWithID(ctx, "[Service: SendAlertReachParticipantsToAdmin] Called")
	alert := entity.Alert{
		Subject: "Question Reached Participants",
		Message: fmt.Sprintf("Question: %s\nTotal Participants: %d people\n1.%s : %d people\n2.%s : %d people", questionText, totalParticipants, firstChoice, firstChoiceCount, secondChoice, secondChoiceCount),
	}

	err := a.sender.SendAdminAlert(ctx, alert)
	if err != nil {
		a.log.ErrorWithID(ctx, "[Service: SendAlertReachParticipantsToAdmin] Failed to send alert:", err)
		return err
	}
	a.log.InfoWithID(ctx, "[Service: SendAlertReachParticipantsToAdmin] Alert sent to SNS!")
	return nil
}

func (a *NotificationService) NotifyUserOfAdminQuestion(ctx context.Context, email, subject, message string) error {
	a.log.InfoWithID(ctx, "[Service: NotifyUserOfAdminQuestion] Called")
	alert := entity.Alert{
		Subject: subject,
		Message: message,
	}
	err := a.sender.SendUserAlert(ctx, alert)
	if err != nil {
		a.log.ErrorWithID(ctx, "[Service: NotifyUserOfAdminQuestion] Failed to send alert:", err)
		return err
	}
	return nil
}

func (a *NotificationService) CheckIsAdmin(ctx context.Context, email string) (bool, error) {
	a.log.InfoWithID(ctx, "[Service: CheckIsAdmin] Called")
	admins, err := a.sender.GetAdminSubscriptions(ctx)
	if err != nil {
		a.log.ErrorWithID(ctx, "[Service: CheckIsAdmin] Failed to get admin subscriptions:", err)
		return false, err
	}

	a.log.InfoWithID(ctx, "[Service: CheckIsAdmin] Found admins:", admins)
	return slices.Contains(admins, email), nil
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
