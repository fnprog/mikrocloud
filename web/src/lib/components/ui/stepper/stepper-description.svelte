<script lang="ts">
  import { cn } from '$lib/utils';
  import { stepItemContext } from './stepper-context.js';

  interface Props {
    class?: string;
    children: any;
  }

  let { class: className, children }: Props = $props();

  let itemCtx = $stepItemContext;
  if (!itemCtx) throw new Error('StepperDescription must be used within a StepperItem');

  let { state } = itemCtx;
</script>

<div data-slot="stepper-description" data-state={state} class={cn('text-sm text-muted-foreground', className)}>
  {@render children()}
</div>