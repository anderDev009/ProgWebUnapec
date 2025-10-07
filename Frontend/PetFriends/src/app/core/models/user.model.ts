export type UserRole = 'adopter' | 'shelter' | 'admin';

export interface User {
  id: number;
  name: string;
  email: string;
  role: UserRole;
  city?: string | null;
  phone?: string | null;
  shelterName?: string | null;
  isApproved: boolean;
}
