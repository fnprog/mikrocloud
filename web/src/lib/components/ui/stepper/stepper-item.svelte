<script lang="ts">
  import { cn } from '$lib/utils';
  import { stepperContext, stepItemContext, type StepState } from './stepper-context.js';

  interface Props {
    step: number;
    completed?: boolean;
    disabled?: boolean;
    loading?: boolean;
    class?: string;
    children: any;
  }

  let {
    step,
    completed = false,
    disabled = false,
    loading = false,
    class: className,
    children,
    ...props
  }: Props = $props();

  let ctx = $stepperContext;
  if (!ctx) throw new Error('StepperItem must be used within a Stepper');

  let { activeStep, registerStep } = ctx;

  let state: StepState = $derived(
    completed || step < activeStep ? 'completed' : activeStep === step ? 'active' : 'inactive'
  );

  let isLoading = $derived(loading && step === activeStep);

  let itemContextValue = $derived({
    step,
    state,
    isDisabled: disabled,
    isLoading,
  });

  // Set context
  $effect(() => {
    stepItemContext.set(itemContextValue);
  });

  // Register this step
  $effect(() => {
    registerStep(step);
  });
</script>

<div
  data-slot="stepper-item"
  class={cn(
    'group/step flex items-center justify-center group-data-[orientation=horizontal]/stepper-nav:flex-row group-data-[orientation=vertical]/stepper-nav:flex-col not-last:flex-1',
    className,
  )}
  data-state={state}
  {...(isLoading ? { 'data-loading': true } : {})}
  {...props}
>
  {@render children()}
</div>