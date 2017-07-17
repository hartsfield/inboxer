/*
* inboxer is a Go package for checking your gmail inbox, it has the following
* features:
*
*  - Mark emails (read/unread/important/etc)
*  - Get labels used in inbox
*  - Get emails by query (eg "in:sent after:2017/01/01 before:2017/01/30")
*  - Get email metadata
*  - Get email main body ("text/plain", "text/html")
*  - Get the number of unread messages
*  - Convert email dates to human readable format
*
*******************************************************************************
*  USE
*******************************************************************************
* CREDENTIALS:
*
* For inboxer to work you must have a gmail account and a file named
* "client_secret.json" containing your authorization info in the root directory
* of your project. To obtain credentials please see step one of this guide:
* https://developers.google.com/gmail/api/quickstart/go
*
*  >Step 1: Turn on the Gmail API

*  >Use this wizard (https://console.developers.google.com/start/api?id=gmail) to create or select a project in the Google Developers Console and automatically turn on the API. Click Continue, then Go to credentials.
*  >On the Add credentials to your project page, click the Cancel button.
*  >At the top of the page, select the OAuth consent screen tab. Select an Email address, enter a Product name if not already set, and click the Save button.
*  >Select the Credentials tab, click the Create credentials button and select OAuth client ID.
*  >Select the application type Other, enter the name "Gmail API Quickstart", and click the Create button.
*  >Click OK to dismiss the resulting dialog.
*  >Click the file_download (Download JSON) button to the right of the client ID.
*  >Move this file to your working directory and rename it client_secret.json.

*
*
* package main
*
* import
*
*******************************************************************************
*  MARKING EMAILS
*******************************************************************************
*
*******************************************************************************
*  GETTING LABELS
*******************************************************************************
*
*******************************************************************************
*  QUERIES
*******************************************************************************
*
*******************************************************************************
*  METADATA
*******************************************************************************
*
*******************************************************************************
*  GETTING THE EMAIL BODY
*******************************************************************************
*
*******************************************************************************
*  GETTING THE NUMBER OF UNREAD MESSAGES
*******************************************************************************
*
*******************************************************************************
*  CONVERTING DATES
*******************************************************************************
*
*******************************************************************************
*  MARK ALL "UNREAD" EMAILS AS "READ"
*******************************************************************************
 */
package inboxer
