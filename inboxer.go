package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.com/hartsfield/gmailAPI"
	gmail "google.golang.org/api/gmail/v1"
)

func main() {
	// Connect to the gmail API service.
	ctx := context.Background()
	srv := gmailAPI.ConnectToService(ctx, gmail.MailGoogleComScope)

	// Get the messages
	msgs, err := srv.Users.Messages.List("me").Do()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("--------------------------------------------------------")
	// Range over the messages
	for _, v := range msgs.Messages[:10] {
		msg, _ := srv.Users.Messages.Get("me", v.Id).Do()
		info := getSender2(msg)
		fmt.Println(info.From, info.Sender, info.Subject)
		// fmt.Println(getSender(msg), getTime(msg.InternalDate))
		if getByLabel(strings.ToUpper("UNREAD"), msg) {
			// fmt.Println(msg.Snippet)
			// fmt.Println("--------------------------------------------------------")
		}
	}
}

func getBody(msg *gmail.Message) (string, error) {
	dec, err := decodeEmailBody(msg.Payload.Parts[0].Body.Data)
	if err != nil {
		return "", err
	}
	return dec, nil
}

func getByLabel(label string, msg *gmail.Message) bool {
	for _, v := range msg.LabelIds {
		if v == label {
			return true
		}
	}
	return false
}

type senderInfo struct {
	Sender, From, To, CC, Subject, MailingList, DeliveredTo, ThreadTopic []string
}

func getSender2(msg *gmail.Message) *senderInfo {
	info := &senderInfo{}
	fmt.Println("--------------------------------------------------------")
	for _, v := range msg.Payload.Headers {
		switch v.Name {
		case "Sender":
			info.Sender = append(info.Sender, v.Value)
		case "From":
			info.From = append(info.From, v.Value)
		case "To":
			info.To = append(info.To, v.Value)
		case "CC":
			info.CC = append(info.CC, v.Value)
		case "Subject":
			info.Subject = append(info.Subject, v.Value)
		case "Mailing-list":
			info.MailingList = append(info.MailingList, v.Value)
		case "Delivered-To":
			info.DeliveredTo = append(info.DeliveredTo, v.Value)
		case "Thread-Topic":
			info.ThreadTopic = append(info.ThreadTopic, v.Value)
		}
	}
	return info
}

func decodeEmailBody(data string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		fmt.Println("decode error:", err)
		return "", err
	}
	return string(decoded), nil
}

func getTime(datetime int64) time.Time {
	conv := strconv.FormatInt(datetime, 10)
	cut := conv[:len(conv)-3]
	tc, err := strconv.ParseInt(cut, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	return time.Unix(tc, 0)
}

// func watchInbox() {
// 	req := &gmail.WatchRequest{
// 		LabelFilterAction: "include",
// 		LabelIds:          []string{"UNREAD"},
// 		TopicName:         "gmailmsg",
// 	}

// 	wr, _ := srv.Users.Watch("me", req).Do()
// 	fmt.Println(wr.ForceSendFields)
// }

// func getMessages() (*gmail.ListMessagesResponse, error) {
// 	// Connect to the gmail API service.
// 	ctx := context.Background()
// 	srv := gmailAPI.ConnectToService(ctx, gmail.MailGoogleComScope)

// 	// Get the messages
// 	msgs, err := srv.Users.Messages.List("me").Do()
// 	if err != nil {
// 		return &gmail.ListMessagesResponse{}, err
// 	}

// 	return msgs, nil
// }
