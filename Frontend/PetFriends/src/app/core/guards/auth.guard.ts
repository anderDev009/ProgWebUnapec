import { inject } from '@angular/core';
import { CanActivateFn, Router, UrlTree } from '@angular/router';

import { UserRole } from '../models/user.model';
import { AuthService } from '../services/auth.service';

export const authGuard: CanActivateFn = (route, state): boolean | UrlTree => {
  const auth = inject(AuthService);
  const router = inject(Router);

  if (!auth.isAuthenticated()) {
    return router.createUrlTree(['/auth/login'], {
      queryParams: { redirectTo: state.url },
    });
  }

  const allowedRoles = route.data?.['roles'] as UserRole[] | undefined;
  const role = auth.currentRole();
  if (allowedRoles?.length && (!role || !allowedRoles.includes(role))) {
    return router.createUrlTree(['/']);
  }
  return true;
};
