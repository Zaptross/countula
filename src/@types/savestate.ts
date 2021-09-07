export type Id = string;

export type Score = {
    highest: number;
    longestStreak: number;
    totalScore: number;
    currentScore: number;
    mistakes: number;
    biggestMistake: number;
};

export interface ScorelessPlayer {
    id: Id;
    username: string;
}

export type Player = ScorelessPlayer & {
    stats: Score;
};

export interface Message {
    id?: Id;
    playerId?: Id;
    text?: string;
}

export interface LastHighScore {
    id?: Id;
    score?: number;
}

export interface State {
    players: Record<Id, Player>;
    lastMessage: Message;
    lastHighScoreReact: LastHighScore;
}

export interface SaveState {
    last: {
        message: {
            messageId: Id | undefined;
            secondMessageId: Id | undefined;
            thirdMessageId: Id | undefined;
        };
        highScoreReact: {
            message: Id | undefined;
        };
    };
    current: {
        number: number;
        secondLast: number;
        thirdLast: number;
        highScore: number;
    };
    players: Record<Id, Player>;
    rules: string[];
}
