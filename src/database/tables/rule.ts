import { BaseEntity, Column, Entity, PrimaryColumn } from 'typeorm';

@Entity()
export class Rule extends BaseEntity {
    @PrimaryColumn({ type: 'int' })
    id: number;

    @Column({ type: 'varchar', length: 50 })
    name: string;

    @Column({ type: 'int' })
    gamesActive: number;

    @Column({ type: 'int' })
    timesVetoed: number;
}
