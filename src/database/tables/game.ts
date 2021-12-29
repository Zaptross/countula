import { BaseEntity, Column, Entity, PrimaryGeneratedColumn } from 'typeorm';

@Entity()
export class Game extends BaseEntity {
    @PrimaryGeneratedColumn()
    id: number;

    @Column({ type: 'bigint' })
    started: number;

    @Column({ type: 'bigint', nullable: true })
    ended: number;

    @Column({ type: 'int' })
    vetoCount: number;
}
