import { CommonModule } from '@angular/common';
import { Component, computed, inject } from '@angular/core';
import { Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';

import { AuthService } from './core/services/auth.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, RouterLink, RouterLinkActive],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App {
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);

  readonly isAuthenticated = this.auth.isAuthenticated;
  readonly currentUser = this.auth.currentUser;
  readonly roleLabel = computed(() => this.currentUser()?.role ?? '');
  readonly currentYear = new Date().getFullYear();

  readonly mainLinks = [
    { label: 'Inicio', path: '/' },
    { label: 'Mascotas', path: '/' },
  ];

  readonly secondaryLinks = computed(() => {
    const role = this.auth.currentRole();
    if (role === 'adopter' || role === 'shelter') {
      return [{ label: 'Solicitudes', path: '/adoption-requests' }];
    }
    if (role === 'admin') {
      return [{ label: 'Panel Admin', path: '/admin/users' }];
    }
    return [];
  });

  logout(): void {
    this.auth.logout();
    void this.router.navigate(['/']);
  }
}
