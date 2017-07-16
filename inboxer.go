// Copyright (c) 2017 J. Hartsfield

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package inboxer is a Go library for checking email using the google Gmail
// API.
package inboxer

// SCOPE:
// TODO:
// channels and go routines
// Mark as read/unread/important/spam
// Get emails by label
// Get emails by query
//
// tests
// Watch inbox
// DOCS
// README.md
// how-to: add client credentials (for readme)
// Get Previews/snippet (put in docs)
//
// WORKS:
// Get emails by date
// Get emails by sender
// Get emails by recipient
// Get emails by subject
// Get emails by mailing-list
// Get emails by thread-topic
// Get labels
// Check for unread messages
// Convert date to human readable format
// Get Body
//
// DONE:
// LICENSE

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
// error if decoding goes wrong. mimeType is used to indicate whether you want
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

// CheckForUnreadByLabel checks for unread mail maching the speciified label.
// NOTE: When checking your inbox for unread messages, it's not uncommon for
// it return thousands of unread messages that you don't know about. To see them
// in gmail, search your mail for "label:unread". For CheckForUnreadByLabel to
// work properly you need to mark all mail as read either through gmail or
// through the MarkAllAsRead() function found in this library.
func CheckForUnreadByLabel(srv *gmail.Service, label string) (int64, error) {
	inbox, err := srv.Users.Labels.Get("me", label).Do()
	if err != nil {
		return -1, err
	}

	// fmt.Println(label.MessagesUnread)
	if inbox.MessagesUnread == 0 && inbox.ThreadsUnread == 0 {
		return 0, nil
	}

	return inbox.MessagesUnread + inbox.ThreadsUnread, nil
}

// CheckForUnread checks for mail labeled "UNREAD".
// NOTE: When checking your inbox for unread messages, it's not uncommon for
// it return thousands of unread messages that you don't know about. To see them
// in gmail, search your mail for "label:unread". For CheckForUnread to
// work properly you need to mark all mail as read either through gmail or
// through the MarkAllAsRead() function found in this library.
func CheckForUnread(srv *gmail.Service) (int64, error) {
	inbox, err := srv.Users.Labels.Get("me", "UNREAD").Do()
	if err != nil {
		return -1, err
	}

	// fmt.Println(label.MessagesUnread)
	if inbox.MessagesUnread == 0 && inbox.ThreadsUnread == 0 {
		return 0, nil
	}

	return inbox.MessagesUnread + inbox.ThreadsUnread, nil
}

// GetLabels gets a list of the labels used in the users inbox.
func GetLabels(srv *gmail.Service) (*gmail.ListLabelsResponse, error) {
	return srv.Users.Labels.List("me").Do()
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

// GetByDate gets and returns emails within the time frame specified.
func GetByDate(srv *gmail.Service /*start time.Time, end time.Time*/) []*gmail.Message {
	inbox, err := srv.Users.Messages.List("me").Q("in:inbox after:2017/01/01 before:2017/01/30").Do()
	if err != nil {
		fmt.Println(err)
	}
	return getById(srv, inbox)
}

// GetPartialMetadata gets some of the useful metadata from the headers.
func GetPartialMetadata(msg *gmail.Message) *PartialMetadata {
	info := &PartialMetadata{}
	fmt.Println("========================================================")
	fmt.Println(msg.Snippet)
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

func getById(srv *gmail.Service, msgs *gmail.ListMessagesResponse) []*gmail.Message {
	var msgSlice []*gmail.Message
	for _, v := range msgs.Messages {
		msg, _ := srv.Users.Messages.Get("me", v.Id).Do()
		msgSlice = append(msgSlice, msg)
	}
	return msgSlice
}

// GetMessages gets and returns gmail messages
func GetMessages(srv *gmail.Service, howMany uint) ([]*gmail.Message, error) {
	var msgSlice []*gmail.Message

	// Get the messages
	msgs, err := srv.Users.Messages.List("me").MaxResults(int64(howMany)).Do()
	if err != nil {
		return msgSlice, err
	}

	return getById(srv, msgs), nil
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
