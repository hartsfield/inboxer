// Package inboxer is a Go library for checking email using the google Gmail
// API.
package inboxer

// SCOPE:
// TODO:
// Check for unread messages
// Mark as read/unread/important/spam
// Get x number of messages
// Get Previews
// Get labels
// Get emails by label
// Get emails by date
// Get emails by sender
// Get emails by recipient
// Get emails by subject
// Get emails by mailing-list
// Get emails by thread-topic
// Watch inbox
// LICENSE
// README.md
// how-to: add client credentials (for readme)
// tests
//
// DONE:
// Get Body
//

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	gmail "google.golang.org/api/gmail/v1"
)

// GetBody gets, decodes, and returns the body of the email. It returns an
// error if decoding goes wrong. mimeType is used to indicate whether you wnat
// the plain text or html encoding ("text/html", "text/plain").
func GetBody(msg *gmail.Message, mimeType string) (string, error) {
	for _, v := range msg.Payload.Parts {
		if v.MimeType == "multipart/alternative" {
			for _, l := range v.Parts {
				if l.MimeType == mimeType && l.Body.Size >= 1 {
					dec, err := decodeEmailBody(l.Body.Data)
					if err != nil {
						return "", err
					}
					return dec, nil
				}
			}
		}
		if v.MimeType == mimeType && v.Body.Size >= 1 {
			dec, err := decodeEmailBody(v.Body.Data)
			if err != nil {
				return "", err
			}
			return dec, nil
		}
	}
	return "", errors.New("Couldn't Read Body")
}

// HasLabel takes a label and an email and checks if that email has that label
func HasLabel(label string, msg *gmail.Message) bool {
	for _, v := range msg.LabelIds {
		if v == strings.ToUpper(label) {
			return true
		}
	}
	return false
}

// PartialMetadata stores email metadata
type PartialMetadata struct {
	Sender, From, To, CC, Subject, MailingList, DeliveredTo, ThreadTopic []string
}

// GetPartialMetadata gets some of the useful metadata from the headers.
func GetPartialMetadata(msg *gmail.Message) *PartialMetadata {
	info := &PartialMetadata{}
	fmt.Println("========================================================")
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

// decodeEmailBody is used to decode the email body by converting from
// URLEncoded base64 to a string.
func decodeEmailBody(data string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		fmt.Println("decode error:", err)
		return "", err
	}
	return string(decoded), nil
}

// ReceivedTime converts parses and converts a unix time stamp into a human
// readable format ().
func ReceivedTime(datetime int64) time.Time {
	conv := strconv.FormatInt(datetime, 10)
	// Remove trailing zeros.
	conv = conv[:len(conv)-3]
	tc, err := strconv.ParseInt(conv, 10, 64)
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
