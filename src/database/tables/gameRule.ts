import {
    BaseEntity,
    Column,
    Entity,
    JoinColumn,
    ManyToOne,
    PrimaryGeneratedColumn,
} from 'typeorm';
import { Game } from './game';
import { Rule } from './rule';

@Entity()
export class GameRule extends BaseEntity {
    @PrimaryGeneratedColumn()
    id: number;

    @ManyToOne((type) => Game)
    @JoinColumn()
    game: Game;

    @ManyToOne((type) => Rule)
    @JoinColumn()
    rule: Rule;

    @Column({ type: 'varchar', length: 50 })
    ruleSettings: string;
}
