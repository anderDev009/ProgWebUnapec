import { Pet } from './pet.model';
import { User } from './user.model';

export type AdoptionStatus = 'pending' | 'approved' | 'rejected';

export interface AdoptionRequest {
  id: number;
  petId: number;
  adopterId: number;
  status: AdoptionStatus;
  message?: string | null;
  createdAt?: string;
  updatedAt?: string;
  pet?: Pet;
  adopter?: Pick<User, 'id' | 'name' | 'email'>;
}
