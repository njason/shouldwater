package main

import (
	"time"

	"github.com/hanzoai/gochimp3"
)

func createAndSendCampaign(apiKey string, templateId uint, listId string) error {
	client := gochimp3.New(apiKey)
	client.Timeout = 60 * time.Second

	createCampaignRequest := gochimp3.CampaignCreationRequest{
		Type: gochimp3.CAMPAIGN_TYPE_REGULAR,
		Recipients: gochimp3.CampaignCreationRecipients{
			ListId: listId,
			SegmentOptions: gochimp3.CampaignCreationSegmentOptions{
				Match: "all",
				Conditions: []string{},
			},
		},
		Settings: gochimp3.CampaignCreationSettings{
			SubjectLine: "It's time to water the trees!",
			Title: "NYC unestablished tree watering alert",
			FromName: "Work for Nature",
			ReplyTo: "noreply@workfornature.org",
			ToName: "NYC Tree Stewards",
			TemplateId: templateId,
		},
	}
	
	createCampaignResponse, err := client.CreateCampaign(&createCampaignRequest)
	if err != nil {
		return err
	}

	if createCampaignResponse == nil {
		return err
	}

	_, err = client.SendCampaign(createCampaignResponse.ID, nil)
	if err != nil {
		return err
	}
	return nil
}
