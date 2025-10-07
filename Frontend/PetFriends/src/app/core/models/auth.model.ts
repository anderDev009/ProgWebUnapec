import { User } from './user.model';

export interface LoginResponse {
  token: string;
  user: User;
}

export interface RegisterPayload {
  name: string;
  email: string;
  password: string;
  role: 'adopter' | 'shelter';
  shelterName?: string | null;
  phone?: string | null;
  city?: string | null;
}

export interface LoginPayload {
  email: string;
  password: string;
}
