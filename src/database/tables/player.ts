import { BaseEntity, Column, Entity, PrimaryGeneratedColumn } from 'typeorm';

@Entity()
export class Player extends BaseEntity {
    @PrimaryGeneratedColumn()
    id: number;

    @Column({ type: 'varchar', length: 40 })
    discordUuid: string;

    @Column({ type: 'varchar', length: 100 })
    username: string;

    @Column({ type: 'bool', nullable: false })
    admin: boolean;
}
