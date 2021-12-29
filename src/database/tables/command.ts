import { BaseEntity, Column, Entity, PrimaryColumn } from 'typeorm';

@Entity()
export class Command extends BaseEntity {
    @PrimaryColumn({ type: 'int' })
    id: number;

    @Column({ type: 'varchar', length: 50 })
    name: string;

    @Column({ type: 'bool', nullable: false })
    requireAdmin: boolean;

    @Column({ type: 'bool', nullable: false })
    enabled: boolean;
}
