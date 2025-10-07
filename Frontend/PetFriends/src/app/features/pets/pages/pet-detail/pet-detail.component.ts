import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, OnDestroy, OnInit, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { Subscription } from 'rxjs';

import { Pet } from '../../../../core/models/pet.model';
import { AdoptionService } from '../../../../core/services/adoption.service';
import { AuthService } from '../../../../core/services/auth.service';
import { PetService } from '../../../../core/services/pet.service';

@Component({
  standalone: true,
  selector: 'app-pet-detail',
  imports: [CommonModule, ReactiveFormsModule, RouterLink],
  templateUrl: './pet-detail.component.html',
  styleUrl: './pet-detail.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class PetDetailComponent implements OnInit, OnDestroy {
  private readonly route = inject(ActivatedRoute);
  private readonly petService = inject(PetService);
  private readonly adoptionService = inject(AdoptionService);
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);

  private subscriptions = new Subscription();

  readonly pet = signal<Pet | null>(null);
  readonly loading = signal<boolean>(true);
  readonly error = signal<string | null>(null);
  readonly success = signal<string | null>(null);

  readonly canRequestAdoption = computed(() => this.auth.isAuthenticated() && this.auth.currentRole() === 'adopter');
  readonly isShelter = computed(() => this.auth.currentRole() === 'shelter');

  readonly requestForm = this.fb.nonNullable.group({
    message: ['', [Validators.maxLength(600)]],
  });

  ngOnInit(): void {
    const sub = this.route.paramMap.subscribe((params) => {
      const id = Number(params.get('id'));
      if (!Number.isFinite(id)) {
        this.error.set('Identificador inválido');
        return;
      }
      this.fetchPet(id);
    });
    this.subscriptions.add(sub);
  }

  ngOnDestroy(): void {
    this.subscriptions.unsubscribe();
  }

  fetchPet(id: number): void {
    this.loading.set(true);
    this.error.set(null);
    this.petService.getById(id).subscribe({
      next: (pet) => {
        this.pet.set(pet);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('No pudimos cargar la información de la mascota.');
        this.loading.set(false);
      },
    });
  }

  submitRequest(): void {
    const pet = this.pet();
    if (!pet) {
      return;
    }

    if (!this.canRequestAdoption()) {
      this.redirectToLogin(pet.id);
      return;
    }

    this.success.set(null);
    this.error.set(null);

    this.adoptionService
      .create(pet.id, {
        message: this.requestForm.controls.message.value || undefined,
      })
      .subscribe({
        next: () => {
          this.success.set('Enviamos tu solicitud al refugio.');
          this.requestForm.reset();
        },
        error: () => this.error.set('No pudimos enviar la solicitud. Intenta nuevamente.'),
      });
  }

  redirectToLogin(petId: number): void {
    void this.router.navigate(['/auth/login'], {
      queryParams: { redirectTo: `/pets/${petId}` },
    });
  }
}
