package main

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	text:= `👋 Hello, I'm @%s. I can upload telegram files to tmpfiles.org a Temporary File Hosting solution.

🌟 <b>Available Commands</b>

/start to start the bot 
/help to know  how to use this bot 
/about to get information about this bot 
	
`
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf(text, b.User.Username), &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}

func about(b *gotgbot.Bot, ctx *ext.Context) error {
	text := `🤖 <b>Bot Source Code:</b>
The source code for this Telegram bot is hosted on GitHub.

This project is not affiliated with tmpfiles.org
Made with ❤️ by BotsGalaxy
`
	inlineButtons := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				gotgbot.InlineKeyboardButton{
					Text: "📂 __Repository__",
					Url:  "https://github.com/botsgalaxy/TmpFiles-telegram-bot",
				},
				gotgbot.InlineKeyboardButton{
					Text: "👤 __Author__",
					Url:  "https://github.com/botsgalaxy/",
				},
			},

			{ 
				gotgbot.InlineKeyboardButton{ 
					Text: "🧑🏻‍💻 __Contact Developer__",
					Url: "https://t.me/primeakash",
				},
			},

			
		},
	}

	_, err := ctx.EffectiveMessage.Reply(b, text, &gotgbot.SendMessageOpts{
		ParseMode:   "html",
		ReplyMarkup: inlineButtons,
	})
	if err != nil {
		return fmt.Errorf("failed to send about message: %w", err)
	}
	return nil
}


func help(b *gotgbot.Bot, ctx *ext.Context) error { 
	text := `<b>
📤 You can use this telegram bot to automate file uploads.
⏰ All uploaded files are automatically deleted after 60 minutes.	
💎 To upload file send the file in this chat </b>

<i>🚀 Max Upload File Size Limit: 100MB

🚫 Please note that certain file extensions are restricted and cannot be uploaded.

🔗 Upon completion of the upload, you will receive a direct download link for your convenience.</i>
`	
	_,err := ctx.EffectiveMessage.Reply(b,text,&gotgbot.SendMessageOpts{ 
		ParseMode: "html",
	})

	if err != nil {
		return fmt.Errorf("failed to send help message: %w", err)
	}
	return nil
}
