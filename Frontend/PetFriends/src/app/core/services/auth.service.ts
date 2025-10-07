import { Injectable, computed, signal } from '@angular/core';
import { Observable, tap } from 'rxjs';

import { LoginPayload, LoginResponse, RegisterPayload } from '../models/auth.model';
import { User, UserRole } from '../models/user.model';
import { ApiService } from './api.service';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly storageKey = 'petmatch.session';
  private readonly token = signal<string | null>(null);
  private readonly userSignal = signal<User | null>(null);

  readonly isAuthenticated = computed(() => !!this.token());
  readonly currentUser = computed(() => this.userSignal());
  readonly currentRole = computed<UserRole | null>(() => this.userSignal()?.role ?? null);

  constructor(private readonly api: ApiService) {
    this.restoreSession();
  }

  register(payload: RegisterPayload): Observable<{ user: User; message: string }> {
    return this.api.post<{ user: User; message: string }>('/auth/register', payload);
  }

  login(payload: LoginPayload): Observable<LoginResponse> {
    return this.api.post<LoginResponse>('/auth/login', payload).pipe(
      tap((response) => {
        this.setSession(response.token, response.user);
      }),
    );
  }

  logout(): void {
    this.clearSession();
  }

  authHeaders(): Record<string, string> | undefined {
    const token = this.token();
    if (!token) {
      return undefined;
    }
    return { Authorization: `Bearer ${token}` };
  }

  refreshCurrentUser(): Observable<{ user: User }> {
    return this.api
      .get<{ user: User }>('/auth/me', {
        headers: this.authHeaders(),
      })
      .pipe(
        tap({
          next: (response) => this.userSignal.set(response.user),
          error: () => this.clearSession(),
        }),
      );
  }

  private setSession(token: string, user: User): void {
    this.token.set(token);
    this.userSignal.set(user);
    if (typeof window !== 'undefined') {
      window.localStorage.setItem(this.storageKey, token);
    }
  }

  private restoreSession(): void {
    if (typeof window === 'undefined') {
      return;
    }
    const storedToken = window.localStorage.getItem(this.storageKey);
    if (!storedToken) {
      return;
    }
    this.token.set(storedToken);
    this.api
      .get<{ user: User }>('/auth/me', {
        headers: this.authHeaders(),
      })
      .subscribe({
        next: (response) => this.userSignal.set(response.user),
        error: () => this.clearSession(),
      });
  }

  private clearSession(): void {
    this.token.set(null);
    this.userSignal.set(null);
    if (typeof window !== 'undefined') {
      window.localStorage.removeItem(this.storageKey);
    }
  }
}
