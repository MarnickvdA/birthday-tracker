package main

import (
	"context"
	"log"
)

func (cfg *apiConfig) scheduleNotifications(ctx context.Context) {
	log.Println("scheduleNotifications:", "deleting past notifications")
	d, err := cfg.DB.DeletePastBirthdayNotifications(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("scheduleNotifications:", d, "entries deleted")

	log.Println("scheduleNotifications:", "scheduling birthday notifications for tomorrow")
	i, err := cfg.DB.ScheduleBirthdayNotifications(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("scheduleNotifications:", i, "notifications scheduled")
}

func (cfg *apiConfig) pushBirthdayNotification(ctx context.Context) (size int, err error) {
	n, err := cfg.DB.GetScheduledBirthdayNotificationsForToday(ctx)
	if err != nil {
		return 0, err
	}

	if len(n) == 0 {
		log.Println("pushBirthdayNotification:", "no notifications scheduled for today, skipping.")
		return 0, nil
	}

	log.Println("pushBirthdayNotification:", len(n), "notifications scheduled for today, let's send the message!")
	if err := cfg.DB.UpdateBirthdayNotificationsStateForToday(ctx, "in_transit"); err != nil {
		return 0, err
	}

	log.Println("pushBirthdayNotification:", "scheduled notifications state updated to 'in_transit'")
	_, err = sendBirthdaySlackMessage(n)

	if err != nil {
		if err := cfg.DB.UpdateBirthdayNotificationsStateForToday(ctx, "error"); err != nil {
			return 0, err
		}

		return 0, err
	}
	log.Println("pushBirthdayNotification:", "Slack message sent succesfully")

	if err := cfg.DB.UpdateBirthdayNotificationsStateForToday(ctx, "sent"); err != nil {
		return 0, err
	}

	log.Println("pushBirthdayNotification:", "scheduled notifications state updated to 'sent'")

	return len(n), nil
}
