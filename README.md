# Botty

Botty is a simple application that enables users to have a conversation with an OpenAI GPT-3 model on Discord.

## Usage

### Running the bot
To use the bot run the ```cmd/main.go``` file with the following flags set:

- ```-discord [token]``` - The token of your Discord bot
- ```-openai [API key]``` - Your OpenAI API key

We also have two more optional flags that we can use:
- ```-engine [engine name]``` - Change the bot's default OpenAI engine (text-davinci-002)
- ```-channel [channel ID]``` - Limit the bot to one channel only

The bot does not reply to DMs, it requires a channel in a server to function.

### Commands
The bot will reply to every message in a channel, unless it is a command.
Commands are prefixed by ```!```.

The following are all the commands the bot recognizes:

- ```!setname [name]``` - Sets the bot's name
- ```!setidentity [context]``` - Sets the context for the conversation, the more in-depth it is the better the model will perform
- ```!adduser [user name]``` - Add a user to the conversation, the bot currently supports up to 4 people in the same conversation
- ```!removeuser [user name]``` - Remove a user from the conversation
- ```!clear``` - Clears every parameter of the bot (name, context, memory, users)
- ```!clearlogs``` - Clears the bot's memory, useful if the conversation is not what you want it to be and want to start again
- ```!clearidentity``` - Clears the context and memory of the bot

## Disclaimer
Botty stores previous messages in its memory and sends them to OpenAI to generate better responses,
while this is a great feature which enables it to remember things, your OpenAI usage will greatly increase
the longer a conversation is going on.
If you are worried about usage and don't mind the conversation being reset you can either use one of the clear
commands that clears the bot's memory or optionally you could restart the app's process.