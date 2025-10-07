import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, OnInit, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule } from '@angular/forms';

import { User } from '../../../../core/models/user.model';
import { AdminService } from '../../../../core/services/admin.service';

@Component({
  standalone: true,
  selector: 'app-admin-users',
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './users.component.html',
  styleUrl: './users.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class UsersComponent implements OnInit {
  private readonly adminService = inject(AdminService);
  private readonly fb = inject(FormBuilder);

  readonly users = signal<User[]>([]);
  readonly loading = signal<boolean>(true);
  readonly error = signal<string | null>(null);
  readonly success = signal<string | null>(null);

  readonly filters = this.fb.nonNullable.group({
    role: [''],
    approved: [''],
  });

  ngOnInit(): void {
    this.fetchUsers();
  }

  fetchUsers(): void {
    this.success.set(null);
    this.loading.set(true);
    this.error.set(null);
    const values = this.filters.getRawValue();
    this.adminService
      .listUsers({
        role: values.role || undefined,
        approved: values.approved === '' ? undefined : values.approved === 'true',
      })
      .subscribe({
        next: (users) => this.users.set(users),
        error: () => this.error.set('No pudimos cargar los usuarios.'),
        complete: () => this.loading.set(false),
      });
  }

  approve(user: User): void {
    this.success.set(null);
    this.error.set(null);

    this.adminService.approveShelter(user.id).subscribe({
      next: (updated) => {
        this.success.set('Refugio aprobado.');
        this.users.update((items) => items.map((item) => (item.id === updated.id ? updated : item)));
      },
      error: () => this.error.set('No pudimos aprobar el refugio.'),
    });
  }
}
