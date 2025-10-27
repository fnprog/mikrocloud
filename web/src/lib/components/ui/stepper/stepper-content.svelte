<script lang="ts">
  import { cn } from '$lib/utils';
  import { stepperContext } from './stepper-context.js';

  interface Props {
    value: number;
    forceMount?: boolean;
    class?: string;
    children: any;
  }

  let { value, forceMount, class: className, children }: Props = $props();

  let ctx = $stepperContext;
  if (!ctx) throw new Error('StepperContent must be used within a Stepper');

  let { activeStep } = ctx;

  let isActive = $derived(value === activeStep);
</script>

{#if forceMount || isActive}
  <div
    data-slot="stepper-content"
    data-state={activeStep}
    class={cn('w-full', className, !isActive && forceMount && 'hidden')}
    hidden={!isActive && forceMount}
  >
    {@render children()}
  </div>
{/if}