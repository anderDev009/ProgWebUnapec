import { Injectable } from '@angular/core';
import { map, Observable } from 'rxjs';

import { AdoptionRequest, AdoptionStatus } from '../models/adoption-request.model';
import { ApiService } from './api.service';
import { AuthService } from './auth.service';

export interface CreateAdoptionPayload {
  message?: string;
}

@Injectable({ providedIn: 'root' })
export class AdoptionService {
  constructor(private readonly api: ApiService, private readonly auth: AuthService) {}

  create(petId: number, payload: CreateAdoptionPayload): Observable<AdoptionRequest> {
    return this.api
      .post<{ request: AdoptionRequest }>(`/pets/${petId}/adoption-requests`, payload, {
        headers: this.auth.authHeaders(),
      })
      .pipe(map((response) => response.request));
  }

  listForAdopter(): Observable<AdoptionRequest[]> {
    return this.api
      .get<{ requests: AdoptionRequest[] }>('/adoption-requests', { headers: this.auth.authHeaders() })
      .pipe(map((response) => response.requests));
  }

  listForShelter(): Observable<AdoptionRequest[]> {
    return this.api
      .get<{ requests: AdoptionRequest[] }>('/adoption-requests', {
        headers: this.auth.authHeaders(),
      })
      .pipe(map((response) => response.requests));
  }

  updateStatus(id: number, status: AdoptionStatus): Observable<AdoptionRequest> {
    return this.api
      .patch<{ request: AdoptionRequest }>(`/adoption-requests/${id}`, { status }, { headers: this.auth.authHeaders() })
      .pipe(map((response) => response.request));
  }
}
