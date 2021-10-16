import {
    APIMessageContentResolvable,
    EmojiIdentifierResolvable,
    Guild,
    Message,
    MessageReaction,
} from 'discord.js';
import GameState from './savestate';
import { EMOJI } from './emoji';
import { toWords } from 'number-to-words';
import { Debouncer } from './utils';
import { join } from 'path';
import { readFile } from 'fs/promises';
import { log } from './logger';
import {
    checkConditionRules,
    newRuleSet,
    describeRules,
    checkPreconditionRules,
} from './rules';

const highScoreDebouncer = new Debouncer(1000);
let messages: { mistake: [{ reply?: string; message?: string }] };

async function getMessages() {
    messages = JSON.parse(
        (await readFile(join(__dirname, 'messages.json'))).toString()
    );
    log('Loaded messages.json');
}
getMessages();

export function processMessage(
    state: GameState,
    message: Message,
    send: (message: APIMessageContentResolvable) => void,
    handleError: (error: Error, message?: Message) => void,
    guild: Guild
) {
    const errorHandler = (error: Error) => handleError(error, message);
    const reply = (
        replyText: APIMessageContentResolvable,
        replyHandler?: (reply: Message) => void
    ) =>
        message
            .reply(replyText)
            .then(replyHandler || (() => undefined))
            .catch(errorHandler);
    const react = (
        emoji: EmojiIdentifierResolvable,
        reactionHandler?: (reaction: MessageReaction) => void
    ) =>
        message
            .react(emoji)
            .then(reactionHandler || (() => undefined))
            .catch(errorHandler);

    try {
        const player = state.getPlayer(message.author.id);

        let startingNumber =
            parseInt(/^-?(\d+)/.exec(message.cleanContent)?.[0] || '') || 0;

        if (
            message.author.id === process.env.DINGO_USER_ID &&
            message.cleanContent.toLowerCase().includes('bork')
        ) {
            startingNumber = state.currentNumber + state.increment;
        }
        if (
            message.cleanContent.includes('skip a few') &&
            state.currentNumber === 2
        ) {
            startingNumber = state.currentNumber + 96;
        }

        if (message.cleanContent.includes('!stats') && player) {
            reply(
                `\nTotal score: ${player?.stats.totalScore}\nCurrent score: ${player.stats.currentScore}\nNumber of Mistakes: ${player.stats.mistakes}\nBiggest Mistake: ${player.stats.biggestMistake}\nLongest Streak: ${player.stats.longestStreak}`
            );
        }
        if (message.cleanContent.includes('!rules')) {
            reply(describeRules());
        }
        if (message.cleanContent.includes('!state')) {
            reply(
                `The current number is ${state.currentNumber
                }, and the high score is ${state.highScore
                } ${state.lastHighScoreReact ? `held by ${getNickname(
                    guild,
                    state.lastHighScoreReact.message.author.id
                )} ` : ''}${EMOJI.HIGH_SCORE}`
            );
        }
        if (message.cleanContent.includes('!list')) {
            reply(
                `Ask nicely. Or else. You can ask for:\n!stats\n!rules\n!state`
            );
        }

        const updates = checkPreconditionRules({
            state,
            message,
            enteredNumber: startingNumber,
        });

        for (const update in updates) {
            if (update === 'startingNumber') {
                startingNumber = updates[update] as number;
            }
        }

        if (startingNumber !== 0) {
            if (
                checkConditionRules({
                    state,
                    message,
                    enteredNumber: startingNumber,
                })
            ) {
                checkSillyMessage(startingNumber, reply);
                react(EMOJI.CHECK);
                state.currentNumber = startingNumber;

                let newHighScore = false;

                if (state.currentNumber > state.highScore) {
                    newHighScore = true;
                    highScoreDebouncer.debounce(() => {
                        state.highScore = state.currentNumber;
                        send(`${toWords(state.currentNumber)}, Ah-hah-hah!`);
                        if (state.lastHighScoreReact) {
                            state.lastHighScoreReact.remove();
                        }

                        react(EMOJI.HIGH_SCORE, (react: MessageReaction) => {
                            (() => state)().lastHighScoreReact = react;
                        });
                    }, []);
                }

                if (state.lastMessage.author.id === message.author.id) {
                    state.currentStreak++;
                } else {
                    state.currentStreak = 0;
                }

                state.updatePlayerStats(
                    {
                        id: message.author.id,
                        username:
                            getNickname(guild, message.author.id) ||
                            message.author.username,
                    },
                    {
                        currentScore: (player?.stats.currentScore || 0) + 1,
                        totalScore: 1 + (player?.stats.totalScore || 0),
                        highest:
                            (player?.stats.highest || 0) < state.currentNumber
                                ? state.currentNumber
                                : undefined,
                        longestStreak:
                            (player?.stats.longestStreak || 0) <
                                state.currentStreak
                                ? state.currentStreak
                                : undefined,
                    }
                );
            } else {
                state.updatePlayerStats(
                    {
                        id: message.author.id,
                        username:
                            getNickname(guild, message.author.id) ||
                            message.author.username,
                    },
                    {
                        biggestMistake:
                            (player?.stats.biggestMistake || 0) <
                                Math.abs(state.currentNumber)
                                ? Math.abs(state.currentNumber)
                                : undefined,
                        mistakes: (player?.stats.mistakes || 0) + 1,
                    }
                );

                react(EMOJI.CROSS);

                const mistake =
                    messages.mistake[
                    Math.round(
                        Math.random() * (messages.mistake.length - 1)
                    )
                    ];

                if (mistake.reply) {
                    reply(formatMessage(mistake.reply, message.author, guild));
                }
                if (mistake.message) {
                    send(formatMessage(mistake.message, message.author, guild));
                }
                startNewGame(state, send);
            }

            state.lastMessage = message;
        }
    } catch (error) {
        errorHandler(error as Error);
    }
}

function startNewGame(
    state: GameState,
    send: (message: APIMessageContentResolvable) => void
) {
    state.nextGame();
    newRuleSet(state, send);
    send(describeRules());
    send(`Starting at: ${state.currentNumber}`);
}

function checkSillyMessage(int: number, reply: (reply: string) => any) {
    switch (int) {
        case 25:
            reply('https://media.giphy.com/media/l1KueikfEPGbTzZcI/giphy.gif');
            break;
        case 42:
            reply('Ah, I see you also know the answer. Where is your towel?');
            break;
        case 69:
            reply('Nice.');
            break;
        case 123:
            reply('How do you know my license number?');
            break;
        case 314:
            reply('https://media.giphy.com/media/l41lUR5urK4IAk3V6/giphy.gif');
            break;
        case 404:
            reply('https://media.giphy.com/media/9J7tdYltWyXIY/giphy.gif');
            break;
        case 420:
            reply('https://media.giphy.com/media/Jnx5ztK49mHJe/giphy.gif');
            break;
        case 666:
            reply('https://media.giphy.com/media/fteNbx39dKs37ZqCDm/giphy.gif');
            break;

        default:
            break;
    }
}

function getNickname(guild: Guild, userId: string) {
    return guild.member(userId)?.displayName;
}

function formatMessage(raw: string, author: Message['author'], guild: Guild) {
    const nickname = getNickname(guild, author.id);
    return raw
        .replace(/{username}/g, nickname || author.username)
        .replace(/{@user}/g, `<@${author.id}>`);
}
