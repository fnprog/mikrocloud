<script lang="ts">
  import { cn } from '$lib/utils';
  import { stepperContext, stepItemContext } from './stepper-context.js';

  interface Props {
    class?: string;
    children: any;
  }

  let { class: className, children }: Props = $props();

  let itemCtx = $stepItemContext;
  if (!itemCtx) throw new Error('StepperIndicator must be used within a StepperItem');

  let stepperCtx = $stepperContext;
  if (!stepperCtx) throw new Error('StepperIndicator must be used within a Stepper');

  let { state, isLoading } = itemCtx;
  let { indicators } = stepperCtx;
</script>

<div
  data-slot="stepper-indicator"
  data-state={state}
  class={cn(
    'relative flex items-center overflow-hidden justify-center size-6 shrink-0 border-background bg-accent text-accent-foreground rounded-full text-xs data-[state=completed]:bg-primary data-[state=completed]:text-primary-foreground data-[state=active]:bg-primary data-[state=active]:text-primary-foreground',
    className,
  )}
>
  <div class="absolute">
    {#if indicators && ((isLoading && indicators.loading) || (state === 'completed' && indicators.completed) || (state === 'active' && indicators.active) || (state === 'inactive' && indicators.inactive))}
      {#if isLoading && indicators.loading}
        {@html indicators.loading}
      {:else if state === 'completed' && indicators.completed}
        {@html indicators.completed}
      {:else if state === 'active' && indicators.active}
        {@html indicators.active}
      {:else if state === 'inactive' && indicators.inactive}
        {@html indicators.inactive}
      {/if}
    {:else}
      {@render children()}
    {/if}
  </div>
</div>