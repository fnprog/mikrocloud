<script lang="ts">
  import { cn } from '$lib/utils';
  import { stepperContext, stepItemContext } from './stepper-context.js';
  import { tick } from 'svelte';

  interface Props {
    class?: string;
    children: any;
    tabIndex?: number;
    asChild?: boolean;
  }

  let {
    class: className,
    children,
    tabIndex,
    asChild = false,
    ...props
  }: Props = $props();

  let stepperCtx = $stepperContext;
  if (!stepperCtx) throw new Error('StepperTrigger must be used within a Stepper');

  let itemCtx = $stepItemContext;
  if (!itemCtx) throw new Error('StepperTrigger must be used within a StepperItem');

  let { setActiveStep, activeStep, registerTrigger, triggerNodes, focusNext, focusPrev, focusFirst, focusLast } = stepperCtx;
  let { step, state, isDisabled, isLoading } = itemCtx;

  let isSelected = $derived(activeStep === step);
  let id = $derived(`stepper-tab-${step}`);
  let panelId = $derived(`stepper-panel-${step}`);

  // Register this trigger for keyboard navigation
  let btnRef: HTMLElement | undefined = $state();

  $effect(() => {
    if (btnRef) {
      registerTrigger(btnRef);
    }
  });

  // Find our index among triggers for navigation
  let myIdx = $derived(triggerNodes.findIndex((n) => n === btnRef));

  function handleKeyDown(e: KeyboardEvent) {
    switch (e.key) {
      case 'ArrowRight':
      case 'ArrowDown':
        e.preventDefault();
        if (myIdx !== -1) focusNext(myIdx);
        break;
      case 'ArrowLeft':
      case 'ArrowUp':
        e.preventDefault();
        if (myIdx !== -1) focusPrev(myIdx);
        break;
      case 'Home':
        e.preventDefault();
        focusFirst();
        break;
      case 'End':
        e.preventDefault();
        focusLast();
        break;
      case 'Enter':
      case ' ':
        e.preventDefault();
        setActiveStep(step);
        break;
    }
  }

  <button
    bind:this={btnRef}
    role="tab"
    {id}
    aria-selected={isSelected}
    aria-controls={panelId}
    tabindex={typeof tabIndex === 'number' ? tabIndex : isSelected ? 0 : -1}
    data-slot="stepper-trigger"
    data-state={state}
    data-loading={isLoading}
    class={cn(
      'cursor-pointer focus-visible:border-ring focus-visible:ring-ring/50 inline-flex items-center gap-3 rounded-full outline-none focus-visible:z-10 focus-visible:ring-[3px] disabled:pointer-events-none disabled:opacity-60',
      className,
    )}
    onclick={() => setActiveStep(step)}
    onkeydown={handleKeyDown}
    disabled={isDisabled}
    {...props}
  >
    {#snippet children({ state, isSelected, isDisabled, isLoading })}
      {@render children()}
    {/snippet}
  </button>
</script>