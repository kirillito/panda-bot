# panda-bot
Slack Bot in Go language

## Usage:
`panda-bot <slack-bot-token>`

## Currently supported commands:
Most commands are matched by substring/regex patterns, so exact message text can fluctuate. Below are working examples
- Set your own vacation: `<bot-name> I'm on vacation from xxxx-xx-xx to xxxx-xx-xx`
- Set user vacation: `<bot-name> <user-tag> is on vacation from xxxx-xx-xx to xxxx-xx-xx`
- Get user vacation data: `<bot-name> list vacations for <user-tag>`
- Goodnight message: `<bot-name> goodnight` (case insensitive)
