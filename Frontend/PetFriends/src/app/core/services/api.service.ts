import { HttpClient, HttpParams } from '@angular/common/http';
import { Inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';

import { API_BASE_URL } from '../config/api.tokens';

@Injectable({ providedIn: 'root' })
export class ApiService {
  constructor(
    private readonly http: HttpClient,
    @Inject(API_BASE_URL) private readonly baseUrl: string,
  ) {}

  get<T>(
    path: string,
    options: {
      params?: Record<string, string | number | boolean | undefined | null>;
      headers?: Record<string, string>;
    } = {},
  ): Observable<T> {
    const httpParams = this.buildParams(options.params);
    return this.http.get<T>(this.resolve(path), {
      params: httpParams,
      headers: options.headers,
    });
  }

  post<T>(path: string, body: unknown, options: { headers?: Record<string, string> } = {}): Observable<T> {
    return this.http.post<T>(this.resolve(path), body, { headers: options.headers });
  }

  put<T>(path: string, body: unknown, options: { headers?: Record<string, string> } = {}): Observable<T> {
    return this.http.put<T>(this.resolve(path), body, { headers: options.headers });
  }

  patch<T>(path: string, body: unknown, options: { headers?: Record<string, string> } = {}): Observable<T> {
    return this.http.patch<T>(this.resolve(path), body, { headers: options.headers });
  }

  delete<T>(path: string, options: { headers?: Record<string, string> } = {}): Observable<T> {
    return this.http.delete<T>(this.resolve(path), { headers: options.headers });
  }

  private resolve(path: string): string {
    if (!path.startsWith('/')) {
      return `${this.baseUrl}/${path}`;
    }
    return `${this.baseUrl}${path}`;
  }

  private buildParams(params?: Record<string, string | number | boolean | undefined | null>): HttpParams | undefined {
    if (!params) {
      return undefined;
    }
    let httpParams = new HttpParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        httpParams = httpParams.set(key, String(value));
      }
    });
    return httpParams;
  }
}
