<script lang="ts">
  import { cn } from '$lib/utils';
  import { stepperContext } from './stepper-context.js';

  interface Props {
    class?: string;
    children: any;
  }

  let { class: className, children }: Props = $props();

  let ctx = $stepperContext;
  if (!ctx) throw new Error('StepperNav must be used within a Stepper');

  let { activeStep, orientation } = ctx;
</script>

<nav
  data-slot="stepper-nav"
  data-state={activeStep}
  data-orientation={orientation}
  class={cn(
    'group/stepper-nav inline-flex data-[orientation=horizontal]:w-full data-[orientation=horizontal]:flex-row data-[orientation=vertical]:flex-col',
    className,
  )}
>
  {@render children()}
</nav>