<script lang="ts">
  import { cn } from '$lib/utils';
  import { stepperContext } from './stepper-context.js';

  interface Props {
    class?: string;
    children: any;
  }

  let { class: className, children }: Props = $props();

  let ctx = $stepperContext;
  if (!ctx) throw new Error('StepperPanel must be used within a Stepper');

  let { activeStep } = ctx;
</script>

<div data-slot="stepper-panel" data-state={activeStep} class={cn('w-full', className)}>
  {@render children()}
</div>