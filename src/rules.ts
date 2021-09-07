import { toWords } from 'number-to-words';
import { deromanize } from 'romans';
import { Rule, CheckParams, GetRule, LoadParams } from './@types/rules';
import GameState from './savestate';

const successConditionRules = (
    state: GameState,
    send: LoadParams['send']
): Rule[] => {
    const active: Rule[] = [];
    if (Math.random() >= 0.3) {
        active.push(rulesByName.takeTurns({ state }));
    }
    if (Math.random() >= 0.85) {
        active.push(rulesByName.romanNumerals({ state }));
    }

    // Randomly pick an increment to count in
    const weighted = [
        1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 3, 3, 4, 6, 7, -1, -3, 112358, 112358,
        112358,
    ];
    const randomNumber = Math.floor(Math.random() * weighted.length);
    active.push(
        randomNumber === 112358
            ? rulesByName.fibonacci({ state })
            : rulesByName.increment({
                  num: weighted[randomNumber],
                  state,
              })
    );

    for (const rule of active) {
        if (rule.onLoad) {
            rule.onLoad({ state, send });
        }
    }

    return active;
};

export function checkConditionRules(checkParams: CheckParams) {
    for (const rule of rulesInPlay.condition) {
        if (rule.check(checkParams) === false) {
            return false;
        }
        return true;
    }
}

export function checkPreconditionRules(checkParams: CheckParams) {
    const updates: Record<string, any> = {};

    for (const rule of rulesInPlay.precondition) {
        if (rule.check(checkParams) && rule.passAction) {
            const result = rule.passAction(checkParams);
            for (const update in result) {
                updates[update] = result[update];
            }
        } else {
            if (rule.failAction) {
                const result = rule.failAction(checkParams);
                for (const update in result) {
                    updates[update] = result[update];
                }
            }
        }
    }

    return updates;
}

export function newRuleSet(state: GameState, send: LoadParams['send']) {
    const newRules = successConditionRules(state, send);
    for (const type in rulesInPlay) {
        rulesInPlay[type as Rule['loadType']] = newRules.filter(
            (rule) => rule.loadType === type
        );
    }
}

export function getActiveRules() {
    return [...rulesInPlay.condition, ...rulesInPlay.precondition];
}

export function getSaveableRules() {
    return getActiveRules().map((rule) => rule.loadName ?? rule.name);
}

export function loadRules(
    ruleNames: string[],
    state: GameState,
    send: LoadParams['send']
) {
    for (const ruleName of ruleNames) {
        let rule: Rule;
        const [name, number] = ruleName.split('::');
        if (rulesByName[name]) {
            if (ruleName.includes('::')) {
                rule = rulesByName[name]({
                    num: parseInt(number),
                    state,
                });
                rulesInPlay[rule.loadType].push(rule);
            } else {
                rule = rulesByName[ruleName]({ state });
                rulesInPlay[rule.loadType].push(rule);
            }

            if (rule.onLoad) {
                rule.onLoad({ state, send });
            }
        }
    }

    let noRules = true;
    for (const type of Object.values(rulesInPlay)) {
        if (type.length > 0) {
            noRules = false;
        }
    }
    // Ensure we are always at least counting in ones
    if (noRules) {
        rulesInPlay.condition.push(rulesByName.increment({ state }));
    }
}

export function removeRules(rules: string[]) {
    for (const type of Object.keys(rulesInPlay)) {
        rulesInPlay[type as Rule['loadType']] = rulesInPlay[
            type as Rule['loadType']
        ].filter((rule) => !rules.includes(rule.loadName!));
    }
}

export function describeRules() {
    const activeRules = getActiveRules();

    const rulesDescriptions = ["Count under Countula's Edicts:"];

    for (const rule of activeRules) {
        rulesDescriptions.push(`***${rule.name}***: ${rule.description}`);
    }

    return rulesDescriptions.join('\n');
}

const rulesByName: Record<string, GetRule> = {
    increment: ({ num, state }): Rule => {
        state.increment = num ?? 1;
        return {
            name: `In ${toWords(state.increment)}s`,
            loadName: `increment::${state.increment}`,
            loadType: 'condition',
            description: `Count in ${toWords(state.increment)}s.`,
            check: ({ state, enteredNumber }) =>
                enteredNumber === state.currentNumber + state.increment,
        };
    },
    takeTurns: () => ({
        name: 'Take Turns',
        loadName: 'takeTurns',
        loadType: 'condition',
        description:
            'You **must** allow another player to take a turn after yours.',
        check: ({ state, message }) =>
            state.lastMessage.author.id !== message.author.id,
    }),
    romanNumerals: () => ({
        name: 'Hail Caesar!',
        loadName: 'romanNumerals',
        loadType: 'precondition',
        description: 'Count using roman numerals.',
        check: ({ state, message }: CheckParams) => {
            const numeral = /^-?([IVXLCDM]+)/.exec(message.cleanContent)?.[1];
            if (
                numeral &&
                deromanize(numeral) === state.currentNumber + state.increment
            ) {
                return true;
            } else {
                return false;
            }
        },
        passAction: ({ state, message }) => {
            const numeral = /^-?([IVXLCDM]+)/.exec(message.cleanContent)?.[1];
            const startingNumber =
                message.cleanContent[0] === '-'
                    ? 0 - deromanize(numeral!)
                    : deromanize(numeral!);
            return { startingNumber };
        },
        failAction: ({ state, message }) => {
            const numeral = /^-?([IVXLCDM]+)/.exec(message.cleanContent)?.[1];
            if (numeral) {
                return { startingNumber: deromanize(numeral) || 0.1 };
            }
        },
    }),
    fibonacci: () => ({
        name: "Fibonacci's Sequence",
        loadName: 'fibonacci',
        loadType: 'condition',
        description: 'Count up by adding the last two numbers together.',
        check: ({ state, enteredNumber }: CheckParams) => {
            if (
                enteredNumber ===
                state.currentNumber + state.secondLastNumber
            ) {
                state.increment = state.secondLastNumber;
                return true;
            } else {
                return false;
            }
        },
        onLoad: ({ state, send }) => {
            if (
                !(
                    state.currentNumber ===
                    state.secondLastNumber + state.thirdLastNumber
                )
            ) {
                // Assigning current numbers in order to create the correct history for the fibonacci sequence
                state.currentNumber = 0;
                state.currentNumber = 1;
                send('The last number was 0');
                send('The current number is 1');
            }
        },
    }),
};

const rulesInPlay: Record<Rule['loadType'], Rule[]> = {
    condition: [],
    precondition: [],
};
