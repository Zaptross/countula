import {
    BaseEntity,
    Column,
    Entity,
    JoinColumn,
    ManyToOne,
    PrimaryGeneratedColumn,
} from 'typeorm';
import { Game } from './game';
import { Player } from './player';

@Entity()
export class Turn extends BaseEntity {
    @PrimaryGeneratedColumn()
    id: number;

    @ManyToOne((type) => Player)
    @JoinColumn()
    player: Player;

    @ManyToOne((type) => Game)
    @JoinColumn()
    game: Game;

    @Column({ type: 'varchar', length: 40 })
    messageUuid: string;

    @Column({ type: 'varchar', length: 100 })
    numberText: string;
}
