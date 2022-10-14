# Karmabot

Karmabot is a Slack App written in Go. In order to run the karmabot, you will first need to create a new Slack App. 

## Slack App Configuration

- Visit https://api.slack.com/apps to create your Slack App.
- Choose "From an app manifest".
- Specify the workspace where you want to add the app.
- Copy and paste the contents of the app-manifest.yml file.
- Review the summary and click "Create".

## Configuring your environment

Karmabot expects the following environment variables

* ENVIRONMENT - This should be either "development" or "production". 
  * This only affects log level. DEBUG and up will be logged for development, INFO and up for production.
* SLACK_BOT_TOKEN - Slack Bot User OAuth token 
  * Found under Features - OAuth & Permissions - Bot User OAuth Token. Will begin with `xoxb`.
* SLACK_BOT_USER - Slack Bot User ID
  * Found by calling the slack auth.test API with your SLACK_BOT_TOKEN (https://api.slack.com/methods/auth.test)
* SLACK_APP_TOKEN - Slack App-Level Token
  * Found under Basic Information - App-Level Tokens. Click on the name of the token. Will begin with `xapp`.
* KARMAFILE - Path to the file where the karma data is stored.

## Testing

From the command line, use `make test` to run all tests.

## Run Karmabot

From the command line, use `make run` to run the app.

## Future work

### Dev

- Maybe write some tests
- Fix custom users to avoid insanity when code is copy/pasted into Slack
- Refactor so we don't have to pass the logger around everywhere

### Features

- Milestones ("You've reached 100/500/1000 Karma, huzzah!)
- Specify your Hogwarts House (for the karma message)
- Leaderboard