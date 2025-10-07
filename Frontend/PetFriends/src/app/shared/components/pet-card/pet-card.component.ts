import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { RouterLink } from '@angular/router';

import { Pet } from '../../../core/models/pet.model';

@Component({
  selector: 'app-pet-card',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './pet-card.component.html',
  styleUrl: './pet-card.component.scss',
})
export class PetCardComponent {
  @Input({ required: true }) pet!: Pet;
  @Input() actionLabel = 'Solicitar adopción';
  @Input() showActions = true;
  @Output() adopt = new EventEmitter<Pet>();

  onAdopt(): void {
    this.adopt.emit(this.pet);
  }
}
