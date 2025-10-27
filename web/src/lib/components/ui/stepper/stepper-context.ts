import { writable } from 'svelte/store';

// Types
export type StepperOrientation = 'horizontal' | 'vertical';
export type StepState = 'active' | 'completed' | 'inactive' | 'loading';
export type StepIndicators = {
  active?: string | HTMLElement;
  completed?: string | HTMLElement;
  inactive?: string | HTMLElement;
  loading?: string | HTMLElement;
};

export interface StepperContextValue {
  activeStep: number;
  setActiveStep: (step: number) => void;
  stepsCount: number;
  orientation: StepperOrientation;
  registerTrigger: (node: HTMLElement | null) => void;
  registerStep: (step: number) => void;
  triggerNodes: HTMLElement[];
  focusNext: (currentIdx: number) => void;
  focusPrev: (currentIdx: number) => void;
  focusFirst: () => void;
  focusLast: () => void;
  indicators: StepIndicators;
}

export interface StepItemContextValue {
  step: number;
  state: StepState;
  isDisabled: boolean;
  isLoading: boolean;
}

// Stores
export const stepperContext = writable<StepperContextValue | undefined>(undefined);
export const stepItemContext = writable<StepItemContextValue | undefined>(undefined);