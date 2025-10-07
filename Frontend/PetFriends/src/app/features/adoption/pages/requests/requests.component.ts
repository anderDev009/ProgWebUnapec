import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, OnInit, signal } from '@angular/core';
import { RouterLink } from '@angular/router';

import { AdoptionRequest, AdoptionStatus } from '../../../../core/models/adoption-request.model';
import { AuthService } from '../../../../core/services/auth.service';
import { AdoptionService } from '../../../../core/services/adoption.service';

@Component({
  standalone: true,
  selector: 'app-requests',
  imports: [CommonModule, RouterLink],
  templateUrl: './requests.component.html',
  styleUrl: './requests.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RequestsComponent implements OnInit {
  private readonly adoptionService = inject(AdoptionService);
  private readonly auth = inject(AuthService);

  readonly requests = signal<AdoptionRequest[]>([]);
  readonly loading = signal<boolean>(true);
  readonly error = signal<string | null>(null);
  readonly success = signal<string | null>(null);

  readonly isAdopter = computed(() => this.auth.currentRole() === 'adopter');
  readonly isShelter = computed(() => this.auth.currentRole() === 'shelter');

  ngOnInit(): void {
    this.fetchRequests();
  }

  fetchRequests(): void {
    this.loading.set(true);
    this.error.set(null);
    const source = this.isShelter() ? this.adoptionService.listForShelter() : this.adoptionService.listForAdopter();

    source.subscribe({
      next: (requests) => this.requests.set(requests),
      error: () => this.error.set('No pudimos cargar las solicitudes.'),
      complete: () => this.loading.set(false),
    });
  }

  updateStatus(request: AdoptionRequest, status: AdoptionStatus): void {
    if (!this.isShelter()) {
      return;
    }

    this.success.set(null);
    this.error.set(null);

    this.adoptionService.updateStatus(request.id, status).subscribe({
      next: (updated) => {
        this.success.set('Actualizamos el estado de la solicitud.');
        this.requests.update((items) => items.map((item) => (item.id === updated.id ? updated : item)));
      },
      error: () => this.error.set('No pudimos actualizar la solicitud.'),
    });
  }
}
