import dotenv from 'dotenv'
import { APIMessageContentResolvable, Client, ClientOptions, EmojiIdentifierResolvable, Intents, MessageReaction, Message, TextChannel, User } from 'discord.js'
import { toWords } from 'number-to-words'
import { Debouncer } from './utils'
import { getErrorHandler } from './errorHandler'
import { EMOJI } from './emoji'
import { join } from 'path'

dotenv.config({
    path: join(__dirname, '..', '.env')
})

let lastInt = 0
let highScore = 0
let lastCrown: MessageReaction

const { COUNTING_CHANNEL_ID } = process.env

const intents = new Intents()
intents.add(Intents.FLAGS.GUILD_MESSAGES, Intents.FLAGS.GUILD_MESSAGE_REACTIONS, Intents.FLAGS.DIRECT_MESSAGES)

const client = new Client({ intents } as ClientOptions)

let send: (sendText: APIMessageContentResolvable, messageHandler?: (msg: Message) => unknown) => void
let handleError: (error: Error, message?: Message) => void
let COUNTULA: User
const channels: Record<string, TextChannel> = {}
const highScoreDebouncer = new Debouncer(1000)
client.once('ready', () => {
    console.log(`${new Date().toLocaleTimeString()}: Logged in as ${client.user!.tag}`)
    COUNTULA = client.user!
    channels.counting = client.channels.cache.get(COUNTING_CHANNEL_ID!) as TextChannel
    send = (sendText: APIMessageContentResolvable, messageHandler?: (msg: Message) => any) => channels.counting.send(sendText).then(messageHandler || (() => undefined))
    handleError = getErrorHandler(send)
})

client.on('message', message => {
    if (message.channel.id === COUNTING_CHANNEL_ID && message.author.id !== COUNTULA.id) {

        const errorHandler = (error: Error) => handleError(error, message)
        const reply = (replyText: APIMessageContentResolvable, replyHandler?: (reply: Message) => void) => message.reply(replyText).then(replyHandler || (() => undefined)).catch(errorHandler)
        const react = (emoji: EmojiIdentifierResolvable, reactionHandler?: (reaction: MessageReaction) => void) => message.react(emoji).then(reactionHandler || (() => undefined)).catch(errorHandler)

        try {
            const startingNumber = /^([0-9-]+)/.exec(message.cleanContent)?.[0]

            if (startingNumber) {
                const int = parseInt(message.cleanContent)
                if (int === lastInt + 1) {
                    react(EMOJI.CHECK)
                    lastInt = int
                    if (lastInt > highScore) {
                        highScoreDebouncer.debounce(() => {
                            highScore = lastInt
                            send(`${toWords(lastInt)}, Ah-hah-hah!`)
                            if (lastCrown) {
                                lastCrown.remove()
                            }

                            react(EMOJI.HIGH_SCORE, (react: MessageReaction) => { lastCrown = react })
                        }, [])
                    }
                } else {
                    lastInt = 0
                    const user = `<@${message.author.id}>`
                    react(EMOJI.CROSS)
                    reply('Fuck sake')
                    send(`You had to go and fuck it didn't you ${user}? Let's start again from the top, shall we?`)
                    send(lastInt)
                }
            }
        } catch (error) {
            errorHandler(error)
        }
    }
})

client.login(process.env.DISCORD_BOT_TOKEN)
