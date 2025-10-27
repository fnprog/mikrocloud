<script lang="ts">
  import { cn } from '$lib/utils';
  import { stepItemContext } from './stepper-context.js';

  interface Props {
    class?: string;
  }

  let { class: className }: Props = $props();

  let itemCtx = $stepItemContext;
  if (!itemCtx) throw new Error('StepperSeparator must be used within a StepperItem');

  let { state } = itemCtx;
</script>

<div
  data-slot="stepper-separator"
  data-state={state}
  class={cn(
    'm-0.5 rounded-full bg-muted group-data-[orientation=vertical]/stepper-nav:h-12 group-data-[orientation=vertical]/stepper-nav:w-0.5 group-data-[orientation=horizontal]/stepper-nav:h-0.5 group-data-[orientation=horizontal]/stepper-nav:flex-1',
    className,
  )}
></div>