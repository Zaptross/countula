import { readFile, writeFile, stat } from 'fs/promises';
import { join } from 'path';
import {
    Client,
    EmojiIdentifierResolvable,
    Message,
    MessageReaction,
    ReactionEmoji,
    TextChannel,
} from 'discord.js';
import { Debouncer } from './utils';
import { log } from './logger';
import {
    State,
    Player,
    Score,
    Id,
    ScorelessPlayer,
    SaveState,
} from './@types/savestate';
import { EMOJI } from './emoji';
import { getActiveRules, getSaveableRules, loadRules } from './rules';
import { LoadParams } from './@types/rules';

export default class GameState {
    private _lastMessage?: Message;

    public get lastMessage() {
        return this._lastMessage as Message;
    }

    public set lastMessage(message: Message) {
        this._thirdLastMessage = this._secondLastMessage;
        this._secondLastMessage = this._lastMessage;
        this._lastMessage = message;
    }

    private _secondLastMessage?: Message;

    public get secondLastMessage() {
        return this._secondLastMessage as Message;
    }
    private _thirdLastMessage?: Message;

    public get thirdLastMessage() {
        return this._thirdLastMessage as Message;
    }

    private _lastHighScoreReact?: MessageReaction;

    public get lastHighScoreReact() {
        return this._lastHighScoreReact as MessageReaction;
    }

    public set lastHighScoreReact(reaction: MessageReaction) {
        this._lastHighScoreReact = reaction;
    }

    private _players: Record<Id, Player> = {};

    public getPlayer(id: Id): Player | undefined {
        return this._players[id];
    }

    private _increment = 1;
    public get increment() {
        return this._increment;
    }
    public set increment(increment: number) {
        this._increment = increment;
    }

    private _currentNumber = 0;

    public get currentNumber() {
        return this._currentNumber;
    }

    public set currentNumber(number: number) {
        this._thirdLastNumber = this.secondLastNumber;
        this._secondLastNumber = this.currentNumber;
        this._currentNumber = number;
        this.saveDebouncer.debounce(this.saveState.bind(this), []);
    }

    private _secondLastNumber = 0;

    public get secondLastNumber() {
        return this._secondLastNumber;
    }

    private _thirdLastNumber = 0;

    public get thirdLastNumber() {
        return this._thirdLastNumber;
    }

    private _highScore = 0;

    public get highScore() {
        return this._highScore;
    }

    public set highScore(score: number) {
        this._highScore = score;
    }

    private _currentStreak = 0;
    public get currentStreak() {
        return this._currentStreak;
    }
    public set currentStreak(streak: number) {
        this._currentStreak = streak;
    }

    private saveDebouncer: Debouncer;
    private saveFilePath: string;
    private client: Client;
    private static instance: GameState;

    public constructor(client: Client) {
        this.client = client;
        this.saveFilePath = join(__dirname, 'state.json');
        this.saveDebouncer = new Debouncer(1000);
    }

    public saveState() {
        const collectedState: SaveState = {
            last: {
                message: {
                    messageId: this._lastMessage?.id,
                    secondMessageId: this._secondLastMessage?.id,
                    thirdMessageId: this._thirdLastMessage?.id,
                },
                highScoreReact: {
                    message: this._lastHighScoreReact?.message.id,
                },
            },
            current: {
                number: this._currentNumber,
                secondLast: this._secondLastNumber,
                thirdLast: this._thirdLastNumber,
                highScore: this._highScore,
            },
            players: this._players,
            rules: getSaveableRules(),
        };
        const savePromise = writeFile(
            this.saveFilePath,
            JSON.stringify(collectedState)
        );
        log(
            `Saved Game State - Current: ${this._currentNumber}, HighScore: ${this._highScore
            }, Players: ${Object.keys(this._players).length
            }, Rules: ${collectedState.rules.join(', ')}`
        );

        return savePromise;
    }

    public async loadState(send: LoadParams['send']) {
        try {
            const loadedState: SaveState = JSON.parse(
                (await readFile(this.saveFilePath)).toString()
            );

            this._currentNumber = loadedState.current.number;
            this._secondLastNumber = loadedState.current.secondLast;
            this._thirdLastNumber = loadedState.current.thirdLast;
            this._highScore = loadedState.current.highScore;
            this._players = loadedState.players;

            const getMessageById = async (id: string) =>
                await (
                    (await this.client.channels.fetch(
                        process.env.COUNTING_CHANNEL_ID!
                    )) as TextChannel
                ).messages.fetch(id);

            loadRules(loadedState.rules, this, send);

            this._lastMessage = loadedState.last.message.messageId
                ? await getMessageById(loadedState.last.message.messageId)
                : undefined;
            this._secondLastMessage = loadedState.last.message.secondMessageId
                ? await getMessageById(loadedState.last.message.secondMessageId)
                : undefined;
            this._thirdLastMessage = loadedState.last.message.thirdMessageId
                ? await getMessageById(loadedState.last.message.thirdMessageId)
                : undefined;

            const lastHighScoreMessage =
                loadedState.last.highScoreReact.message || true
                    ? await getMessageById(
                        loadedState.last.highScoreReact.message ||
                        '878525356976529408'
                    )
                    : undefined;

            this._lastHighScoreReact = lastHighScoreMessage
                ? lastHighScoreMessage.reactions.cache.find(
                    (x) => x.emoji.name === EMOJI.HIGH_SCORE
                )
                : undefined;

            log(
                `Loaded Game State - Current: ${this._currentNumber
                }, HighScore: ${this._highScore}, Players: ${Object.keys(this._players).length
                }, Rules: ${loadedState.rules.join(', ')}`
            );
        } catch (error) {
            log(error as Error);
        }
    }

    public updatePlayerStats(
        playerData: ScorelessPlayer,
        updates: Partial<Score>
    ) {
        const { id } = playerData;
        if (!this.getPlayer(id)) {
            this.createNewPlayer(playerData);
        }
        const player = this.getPlayer(id);
        if (player) {
            for (const key in updates) {
                const stat = key as keyof Score;
                if (
                    updates[stat] !== undefined &&
                    updates[stat] !== null &&
                    updates[stat] !== player.stats[stat]
                ) {
                    player.stats[stat] = updates[stat] as number;
                }
            }
        }
    }

    private createNewPlayer(author: ScorelessPlayer) {
        if (!this._players[author.id]) {
            const player: Player = {
                id: author.id,
                username: author.username,
                stats: {
                    highest: 0,
                    longestStreak: 0,
                    currentScore: 0,
                    totalScore: 0,
                    mistakes: 0,
                    biggestMistake: 0,
                },
            };

            this._players[author.id] = player;

            return player;
        }
    }

    public nextGame() {
        this.currentNumber = 0;

        // Reset all player's current score to 0
        for (const player of Object.values(this._players)) {
            const { id, username } = player;
            this.updatePlayerStats({ id, username }, { currentScore: 0 });
        }
    }
}
