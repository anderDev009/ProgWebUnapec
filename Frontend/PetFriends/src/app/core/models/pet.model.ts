import { User } from './user.model';

export type PetStatus = 'available' | 'adopted';

export interface Pet {
  id: number;
  name: string;
  species: string;
  breed?: string | null;
  age: number;
  description?: string | null;
  location?: string | null;
  photoURL?: string | null;
  status: PetStatus;
  shelterId: number;
  shelter?: Pick<User, 'id' | 'name' | 'shelterName' | 'city'>;
  createdAt?: string;
}
