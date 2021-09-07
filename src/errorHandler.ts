import { APIMessageContentResolvable, Message, MessageResolvable } from 'discord.js'
import { generateUuid } from './utils'
import { log } from './logger'

export interface ErrorData {
    code: string,
    stack: string,
    error: Error,
    message?: {
        text: string
    }
}

export function getErrorHandler(send: (message: APIMessageContentResolvable) => void) {
    return (error: Error, message?: Message) => {
        const data = {
            code: `E${generateUuid()}`,
            stack: error.stack,
            error,
            message: {
                text: message?.cleanContent
            }
        }
        send(`<@${message ? message.author.id : process.env.ADMINISTRATOR_USER_ID}> I don't know how, but ya fucked it${message && message.author.id !== process.env.ADMINISTRATOR_USER_ID ? `. Tell <@${process.env.ADMINISTRATOR_USER_ID}> to fix his shit:` : ':'} ${data.code}`)
        return log(data as ErrorData)
    }
}
