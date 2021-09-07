import dotenv from 'dotenv';
import {
    APIMessageContentResolvable,
    Client,
    ClientOptions,
    EmojiIdentifierResolvable,
    Intents,
    MessageReaction,
    Message,
    TextChannel,
    User,
    Guild,
} from 'discord.js';
import { getErrorHandler } from './errorHandler';
import { join } from 'path';
import { log } from './logger';
import GameState from './savestate';
import { adminControls } from './adminControls';
import { processMessage } from './gameLogic';
import { EMOJI } from './emoji';

dotenv.config({
    path: join(__dirname, '..', '.env'),
});

const { COUNTING_CHANNEL_ID } = process.env;

const intents = new Intents();
intents.add(
    Intents.FLAGS.GUILD_MESSAGES,
    Intents.FLAGS.GUILD_MESSAGE_REACTIONS,
    Intents.FLAGS.DIRECT_MESSAGES
);

const client = new Client({ intents } as ClientOptions);
const state = new GameState(client);

let guild: Guild;
async function getGuild() {
    if (guild) {
        return guild;
    } else {
        guild = await client.guilds.fetch(process.env.COUNTING_GUILD_ID!);
    }
}
let send: (
    sendText: APIMessageContentResolvable,
    messageHandler?: (msg: Message) => unknown
) => void;
let sendAdmin: (
    sendText: APIMessageContentResolvable,
    messageHandler?: (msg: Message) => unknown
) => void;

let handleError: (error: Error, message?: Message) => void;
let lastError: string | undefined;
let COUNTULA: User;

const channels: Record<string, TextChannel> = {};
client.once('ready', () => {
    log(`${new Date().toLocaleTimeString()}: Logged in as ${client.user!.tag}`);

    COUNTULA = client.user!;

    channels.counting = client.channels.cache.get(
        COUNTING_CHANNEL_ID!
    ) as TextChannel;

    channels.admin = client.channels.cache.get(
        process.env.ADMINISTRATOR_CHANNEL_ID!
    ) as TextChannel;

    send = (
        sendText: APIMessageContentResolvable,
        messageHandler?: (msg: Message) => any
    ) =>
        channels.counting
            .send(sendText)
            .then(messageHandler || (() => undefined));

    sendAdmin = (
        sendText: APIMessageContentResolvable,
        messageHandler?: (msg: Message) => any
    ) =>
        channels.admin.send(sendText).then(messageHandler || (() => undefined));
    const errorHandlerLogFunction = getErrorHandler(send);
    handleError = (error: Error, message?: Message) => {
        lastError = errorHandlerLogFunction(error, message);
    };
    getGuild();
    state.loadState(send);

    send("Ah-hah-hah! It's time again for counting!");
});

client.login(process.env.DISCORD_BOT_TOKEN);

client.on('message', (message) => {
    if (
        message.channel.id === process.env.ADMINISTRATOR_CHANNEL_ID &&
        message.author.id === process.env.ADMINISTRATOR_USER_ID
    ) {
        adminControls(state, sendAdmin, send, message, lastError);
    }

    if (
        message.channel.id === COUNTING_CHANNEL_ID &&
        message.author.id !== COUNTULA.id
    ) {
        processMessage(state, message, send, handleError, guild);
    }
});

function getExitHandler(signal: string) {
    return async () => {
        send(`***yawn*** Time for counting sheep! ${EMOJI.SHEEP}`);
        log(`Shutting down due to ${signal}`);
        await state.saveState();

        client.destroy();
        process.exit();
    };
}

function registerExitHandlers() {
    const exitCases = ['SIGINT', 'SIGUSR1', 'SIGUSR2', 'SIGTERM'];
    for (const exitCode of exitCases) {
        process.on(exitCode, getExitHandler(exitCode));
    }
}
registerExitHandlers();

function setTerminalTitle(title: string) {
    process.stdout.write(
        String.fromCharCode(27) + ']0;' + title + String.fromCharCode(7)
    );
}

setTerminalTitle('COUNTULA');
