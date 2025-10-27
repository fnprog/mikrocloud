<script lang="ts">
  import { cn } from '$lib/utils';
  import { stepperContext, type StepperOrientation, type StepIndicators } from './stepper-context.js';
  import { tick } from 'svelte';

  interface Props {
    defaultValue?: number;
    value?: number;
    onValueChange?: (value: number) => void;
    orientation?: StepperOrientation;
    indicators?: StepIndicators;
    class?: string;
    children: any;
  }

  let {
    defaultValue = 1,
    value,
    onValueChange,
    orientation = 'horizontal',
    class: className,
    children,
    indicators = {},
    ...props
  }: Props = $props();

  let activeStep = $state(defaultValue);
  let triggerNodes = $state<HTMLElement[]>([]);
  let registeredSteps = $state<Set<number>>(new Set());

  // Register/unregister triggers
  function registerTrigger(node: HTMLElement | null) {
    if (node && !triggerNodes.includes(node)) {
      triggerNodes = [...triggerNodes, node];
    } else if (!node) {
      // Remove null nodes
      triggerNodes = triggerNodes.filter(n => n !== null);
    }
  }

  function registerStep(step: number) {
    registeredSteps.add(step);
    // Trigger reactivity
    registeredSteps = new Set(registeredSteps);
  }

  function handleSetActiveStep(step: number) {
    if (value === undefined) {
      activeStep = step;
    }
    onValueChange?.(step);
  }

  let currentStep = $derived(value ?? activeStep);

  // Keyboard navigation logic
  function focusTrigger(idx: number) {
    if (triggerNodes[idx]) triggerNodes[idx].focus();
  }

  function focusNext(currentIdx: number) {
    focusTrigger((currentIdx + 1) % triggerNodes.length);
  }

  function focusPrev(currentIdx: number) {
    focusTrigger((currentIdx - 1 + triggerNodes.length) % triggerNodes.length);
  }

  function focusFirst() {
    focusTrigger(0);
  }

  function focusLast() {
    focusTrigger(triggerNodes.length - 1);
  }

  // Count steps
  let stepsCount = $derived(Math.max(...Array.from(registeredSteps), 0));

  // Context value
  let contextValue = $derived({
    activeStep: currentStep,
    setActiveStep: handleSetActiveStep,
    stepsCount,
    orientation,
    registerTrigger,
    registerStep,
    focusNext,
    focusPrev,
    focusFirst,
    focusLast,
    triggerNodes,
    indicators,
  });

  // Set context
  $effect(() => {
    stepperContext.set(contextValue);
  });
</script>

<div
  role="tablist"
  aria-orientation={orientation}
  data-slot="stepper"
  class={cn('w-full', className)}
  data-orientation={orientation}
  {...props}
>
  {@render children()}
</div>