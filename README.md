# Whisper-bot

> Use Whisper to transcribe

Go CLI to fuels a Telegram bot that lets you send voice notes to be transcribed using Whisper.

## Installation
After you cloned the repo, open the `env.example` file with a text editor and fill in your credentials. 
- `TELEGRAM_TOKEN`: Your Telegram Bot token
  - Follow [this guide](https://core.telegram.org/bots/tutorial#obtain-your-bot-token) to create a bot and get the token.
- `TELEGRAM_ID` (Optional): Your Telegram User ID
  - If you set this, only you will be able to interact with the bot.
  - To get your ID, message `@userinfobot` on Telegram.
- `EDIT_WAIT_SECONDS` (Optional): Amount of seconds to wait between edits
  - This is set to `1` by default, but you can increase if you start getting a lot of `Too Many Requests` errors.
- Save the file, and rename it to `.env`.
> **Note** Make sure you rename the file to _exactly_ `.env`! The program won't work otherwise.

Finally, open the terminal in your computer (if you're on windows, look for `PowerShell`), navigate to the directory and run `./whisper-telegram`.

## License

This repository is licensed under the [MIT License](LICENSE).
