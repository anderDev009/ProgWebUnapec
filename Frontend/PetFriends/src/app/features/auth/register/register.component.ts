import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';

import { RegisterPayload } from '../../../core/models/auth.model';
import { AuthService } from '../../../core/services/auth.service';

@Component({
  standalone: true,
  selector: 'app-register',
  imports: [CommonModule, ReactiveFormsModule, RouterLink],
  templateUrl: './register.component.html',
  styleUrl: './register.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RegisterComponent {
  private readonly fb = inject(FormBuilder);
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);

  readonly loading = signal<boolean>(false);
  readonly error = signal<string | null>(null);
  readonly success = signal<string | null>(null);

  readonly form = this.fb.nonNullable.group({
    name: ['', [Validators.required, Validators.minLength(3)]],
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]],
    role: ['adopter', [Validators.required]],
    shelterName: [''],
    phone: [''],
    city: [''],
  });

  readonly isShelter = computed(() => this.form.controls.role.value === 'shelter');

  submit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    if (this.isShelter() && !this.form.controls.shelterName.value) {
      this.error.set('El nombre del refugio es obligatorio.');
      return;
    }

    this.loading.set(true);
    this.error.set(null);
    this.success.set(null);

    const raw = this.form.getRawValue();
    const role = raw.role as RegisterPayload['role'];

    const payload: RegisterPayload = {
      name: raw.name,
      email: raw.email,
      password: raw.password,
      role,
      shelterName: role === 'shelter' ? raw.shelterName || undefined : undefined,
      phone: raw.phone || undefined,
      city: raw.city || undefined,
    };

    this.auth
      .register(payload)
      .subscribe({
        next: (response) => {
          this.success.set(response.message || 'Registro completado. Puedes iniciar sesion.');
          if (role === 'adopter') {
            void this.router.navigate(['/auth/login'], { queryParams: { email: raw.email } });
          }
        },
        error: (err) => {
          const message = err?.error?.error || 'No pudimos completar el registro.';
          this.error.set(message);
          this.loading.set(false);
        },
        complete: () => this.loading.set(false),
      });
  }
}
