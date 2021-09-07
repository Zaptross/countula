import { Message } from 'discord.js';
import GameState from '../savestate';

export interface CheckParams {
    state: GameState;
    message: Message;
    enteredNumber: number;
}

export interface LoadParams {
    state: GameState;
    send: (message: string) => void;
}

export interface Rule {
    name: string;
    loadName?: string;
    loadType: 'condition' | 'precondition';
    description: string;
    check: (checkParams: CheckParams) => boolean;
    onLoad?: (checkParams: LoadParams) => any;
    passAction?: (checkParams: CheckParams) => any;
    failAction?: (checkParams: CheckParams) => any;
}

export type GetRule = (params: { num?: number; state: GameState }) => Rule;
