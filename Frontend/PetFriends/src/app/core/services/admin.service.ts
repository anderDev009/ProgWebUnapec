import { Injectable } from '@angular/core';
import { map, Observable } from 'rxjs';

import { User } from '../models/user.model';
import { ApiService } from './api.service';
import { AuthService } from './auth.service';

@Injectable({ providedIn: 'root' })
export class AdminService {
  constructor(private readonly api: ApiService, private readonly auth: AuthService) {}

  listUsers(params: { role?: string; approved?: boolean } = {}): Observable<User[]> {
    return this.api
      .get<{ users: User[] }>('/admin/users', {
        params,
        headers: this.auth.authHeaders(),
      })
      .pipe(map((response) => response.users));
  }

  approveShelter(id: number): Observable<User> {
    return this.api
      .post<{ user: User; message: string }>(`/admin/shelters/${id}/approve`, {}, { headers: this.auth.authHeaders() })
      .pipe(map((response) => response.user));
  }
}
