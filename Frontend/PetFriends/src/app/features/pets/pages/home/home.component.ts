import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, OnInit, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';

import { Pet } from '../../../../core/models/pet.model';
import { AuthService } from '../../../../core/services/auth.service';
import { PetFilters, PetService } from '../../../../core/services/pet.service';
import { PetCardComponent } from '../../../../shared/components/pet-card/pet-card.component';

@Component({
  standalone: true,
  selector: 'app-home',
  imports: [CommonModule, ReactiveFormsModule, PetCardComponent],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class HomeComponent implements OnInit {
  private readonly petService = inject(PetService);
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);

  readonly pets = signal<Pet[]>([]);
  readonly loading = signal<boolean>(true);
  readonly error = signal<string | null>(null);

  readonly filterForm = this.fb.nonNullable.group({
    species: [''],
    location: [''],
    minAge: this.fb.control<number | null>(null),
    maxAge: this.fb.control<number | null>(null),
  });

  ngOnInit(): void {
    this.loadPets();
  }

  loadPets(filters: PetFilters = {}): void {
    this.loading.set(true);
    this.error.set(null);
    this.petService.list(filters).subscribe({
      next: (pets) => this.pets.set(pets),
      error: () => this.error.set('No pudimos cargar las mascotas. Inténtalo más tarde.'),
      complete: () => this.loading.set(false),
    });
  }

  onSearch(): void {
    const values = this.filterForm.getRawValue();
    this.loadPets({
      species: values.species || undefined,
      location: values.location || undefined,
      minAge: values.minAge ?? undefined,
      maxAge: values.maxAge ?? undefined,
    });
  }

  resetFilters(): void {
    this.filterForm.reset({
      species: '',
      location: '',
      minAge: null,
      maxAge: null,
    });
    this.loadPets();
  }

  onAdopt(pet: Pet): void {
    if (!this.authService.isAuthenticated()) {
      void this.router.navigate(['/auth/login'], {
        queryParams: { redirectTo: `/pets/${pet.id}` },
      });
      return;
    }
    void this.router.navigate(['/pets', pet.id], {
      queryParams: { focus: 'adoption-form' },
    });
  }
}
