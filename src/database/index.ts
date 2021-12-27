import { createConnection } from 'typeorm';
import { Command } from './tables/command';
import { CommandLog } from './tables/commandLog';
import { Game } from './tables/game';
import { GameRule } from './tables/gameRule';
import { Player } from './tables/player';
import { PlayerStats } from './tables/playerStats';
import { Rule } from './tables/rule';
import { Turn } from './tables/turn';

export async function connect() {
    return createConnection({
        type: 'postgres',
        host: process.env.TYPEORM_HOST!,
        port: Number(process.env.TYPEORM_PORT!),
        username: process.env.TYPEORM_USERNAME!,
        password: process.env.TYPEORM_PASSWORD!,
        database: process.env.TYPEORM_DATABASE!,
        synchronize: true,
        logging: true,
        entities: [
            Command,
            CommandLog,
            Game,
            GameRule,
            Player,
            PlayerStats,
            Rule,
            Turn,
        ],
    });
}
