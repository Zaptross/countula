import {
    BaseEntity,
    Column,
    Entity,
    JoinColumn,
    ManyToOne,
    PrimaryGeneratedColumn,
} from 'typeorm';
import { Command } from './command';
import { Game } from './game';
import { Player } from './player';

@Entity()
export class CommandLog extends BaseEntity {
    @PrimaryGeneratedColumn()
    id: number;

    @ManyToOne((type) => Command)
    @JoinColumn()
    command: Command;

    @ManyToOne((type) => Player)
    @JoinColumn()
    player: Player;

    @ManyToOne((type) => Game)
    @JoinColumn()
    game: Game;

    @Column({ type: 'varchar', length: 40 })
    messageUuid: string;

    @Column({ type: 'varchar', length: 100 })
    messageText: string;
}
