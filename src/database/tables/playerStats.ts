import {
    BaseEntity,
    Column,
    Entity,
    JoinColumn,
    OneToOne,
    PrimaryGeneratedColumn,
} from 'typeorm';
import { Player } from './player';

@Entity()
export class PlayerStats extends BaseEntity {
    @PrimaryGeneratedColumn()
    id: number;

    @OneToOne((type) => Player)
    @JoinColumn()
    player: Player;

    @Column({ type: 'int' })
    gamesPlayed: number;

    @Column({ type: 'int' })
    totalTurns: number;

    @Column({ type: 'int' })
    totalFails: number;

    @Column({ type: 'int' })
    longestStreak: number;

    @Column({ type: 'int' })
    longestStreakBreak: number;
}
