import { Routes } from '@angular/router';

import { authGuard } from './core/guards/auth.guard';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () => import('./features/pets/pages/home/home.component').then((m) => m.HomeComponent),
  },
  {
    path: 'pets/:id',
    loadComponent: () => import('./features/pets/pages/pet-detail/pet-detail.component').then((m) => m.PetDetailComponent),
  },
  {
    path: 'auth/login',
    loadComponent: () => import('./features/auth/login/login.component').then((m) => m.LoginComponent),
  },
  {
    path: 'auth/register',
    loadComponent: () => import('./features/auth/register/register.component').then((m) => m.RegisterComponent),
  },
  {
    path: 'adoption-requests',
    canActivate: [authGuard],
    data: { roles: ['adopter', 'shelter'] },
    loadComponent: () =>
      import('./features/adoption/pages/requests/requests.component').then((m) => m.RequestsComponent),
  },
  {
    path: 'admin/users',
    canActivate: [authGuard],
    data: { roles: ['admin'] },
    loadComponent: () => import('./features/admin/pages/users/users.component').then((m) => m.UsersComponent),
  },
  { path: '**', redirectTo: '' },
];
