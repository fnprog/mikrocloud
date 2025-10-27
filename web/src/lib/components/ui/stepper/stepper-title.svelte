<script lang="ts">
  import { cn } from '$lib/utils';
  import { stepItemContext } from './stepper-context.js';

  interface Props {
    class?: string;
    children: any;
  }

  let { class: className, children }: Props = $props();

  let itemCtx = $stepItemContext;
  if (!itemCtx) throw new Error('StepperTitle must be used within a StepperItem');

  let { state } = itemCtx;
</script>

<h3 data-slot="stepper-title" data-state={state} class={cn('text-sm font-medium leading-none', className)}>
  {@render children()}
</h3>