package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"gitlab.com/hartsfield/gmailAPI"
	gmail "google.golang.org/api/gmail/v1"
)

func main() {
	check()
}

func check() {
	// Connect to the gmail API service.
	ctx := context.Background()
	srv := gmailAPI.ConnectToService(ctx, gmail.MailGoogleComScope)

	// Get the messages
	msgs, err := srv.Users.Messages.List("me").Do()
	if err != nil {
		fmt.Println(err)
	}

	// Range over the messages
	for _, v := range msgs.Messages[:5] {
		msg, _ := srv.Users.Messages.Get("me", v.Id).Do()

		fmt.Println(msg.Snippet)
		fmt.Println(msg.LabelIds)

		// Time
		fmt.Println(getTime(msg.InternalDate))

		// dec, err := decodeEmailBody(msg.Payload.Parts[0].Body.Data)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println(dec)

	}

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
	cut := s[:len(s)-3]
	tc, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	return time.Unix(tc, 0)
}
