import { connect } from '../src/database';
import { Player } from '../src/database/tables/player';
import dotenv from 'dotenv';
import { join } from 'path';
import { Game } from '../src/database/tables/game';
import { GameRule } from '../src/database/tables/gameRule';

dotenv.config({
    path: join(__dirname, '..', '.env'),
});

async function main() {
    await connect();

    await Player.insert({
        discordUuid: '306577431932698625',
        username: 'Zaptross',
        admin: true,
    });

    console.log(
        (await Player.find({ admin: true })).forEach((p) => p.username)
    );

    console.log('');

    await Game.insert({ started: Date.now() });

    console.log(await GameRule.count());
}

main();
