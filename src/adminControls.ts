import { Message } from 'discord.js';
import GameState from './savestate';
import { log } from './logger';
import { describeRules, loadRules, removeRules } from './rules';
import { toWords } from 'number-to-words';

export function adminControls(
    state: GameState,
    sendAdmin: (message: string) => void,
    sendCountingChannel: (message: string) => void,
    message: Message,
    lastError?: string
) {
    log({
        text: message.cleanContent,
        user: message.author.username,
        userId: message.author.id,
    });

    const [command, ...args] = message.cleanContent.split(' ');

    if (command.includes('!list')) {
        message.reply(
            "You forgetting again, boss? There's: \n!setCurrent <number>\n!setHighScore <number>\n!lastError\n!addRule <ruleNames>\n!removeRule <ruleNames>\n!setIncrement <number>"
        );
    }

    if (command.includes('!setCurrent')) {
        state.currentNumber = parseInt(args[0]);
        message.reply(
            `You got it boss. The current number is ${state.currentNumber}`
        );
    }

    if (command.includes('!setHighScore')) {
        state.highScore = parseInt(args[0]);
        message.reply(
            `Changing history again, eh boss? The high score is ${state.highScore}`
        );
    }

    if (command.includes('!lastError')) {
        if (lastError) {
            sendAdmin(lastError);
        } else {
            message.reply('Actually, coast is looking clear boss.');
        }
    }

    if (command.includes('!addRule')) {
        loadRules(args, state, sendCountingChannel);
        message.reply(describeRules());
        state.saveState();
    }

    if (command.includes('!removeRule')) {
        removeRules(args);
        message.reply(describeRules());
        state.saveState();
    }

    if (command.includes('!setIncrement')) {
        state.increment = parseInt(args[0]);
        message.reply(
            `Alrighty boss, counting in ${toWords(state.increment)}s.`
        );
    }
}
