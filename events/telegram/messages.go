package telegram

const MsgHelp = `I can save and keep your pages. Also i can offer you them to read.

In order to save page, just send me a link to it.

In order to get a random page, just send me a command /rnd
Caution! After read your page will be removed`

const MsgHello = "Hi there! \n\n" + MsgHelp

const (
	UnknownCmd      = "Unknown Command 😬"
	MsgNoSavedPages = "You have no saved pages 🙈"
	MsgSaved        = "Saved! 👌🏼"
	MsgAlreadyExist = "You already save have this page in the list 🤓"
)
