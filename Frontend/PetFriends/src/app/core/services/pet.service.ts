import { Injectable } from '@angular/core';
import { map, Observable } from 'rxjs';

import { Pet, PetStatus } from '../models/pet.model';
import { ApiService } from './api.service';
import { AuthService } from './auth.service';

export interface PetFilters {
  species?: string;
  breed?: string;
  location?: string;
  minAge?: number;
  maxAge?: number;
  status?: PetStatus;
}

export interface PetPayload {
  name: string;
  species: string;
  breed?: string;
  age: number;
  description?: string;
  location?: string;
  photoUrl?: string | null;
  status?: PetStatus;
}

@Injectable({ providedIn: 'root' })
export class PetService {
  constructor(private readonly api: ApiService, private readonly auth: AuthService) {}

  list(filters: PetFilters = {}): Observable<Pet[]> {
    const params: Record<string, string | number> = {};
    if (filters.species) {
      params['species'] = filters.species;
    }
    if (filters.breed) {
      params['breed'] = filters.breed;
    }
    if (filters.location) {
      params['location'] = filters.location;
    }
    if (typeof filters.minAge === 'number') {
      params['minAge'] = filters.minAge;
    }
    if (typeof filters.maxAge === 'number') {
      params['maxAge'] = filters.maxAge;
    }
    if (filters.status) {
      params['status'] = filters.status;
    }
    return this.api.get<{ pets: Pet[] }>('/pets', { params }).pipe(map((response) => response.pets));
  }

  getById(id: number): Observable<Pet> {
    return this.api.get<{ pet: Pet }>(`/pets/${id}`).pipe(map((response) => response.pet));
  }

  create(payload: PetPayload): Observable<Pet> {
    return this.api
      .post<{ pet: Pet }>('/pets', payload, { headers: this.auth.authHeaders() })
      .pipe(map((response) => response.pet));
  }

  update(id: number, payload: PetPayload): Observable<Pet> {
    return this.api
      .put<{ pet: Pet }>(`/pets/${id}`, payload, { headers: this.auth.authHeaders() })
      .pipe(map((response) => response.pet));
  }

  remove(id: number): Observable<void> {
    return this.api.delete<void>(`/pets/${id}`, { headers: this.auth.authHeaders() });
  }
}

