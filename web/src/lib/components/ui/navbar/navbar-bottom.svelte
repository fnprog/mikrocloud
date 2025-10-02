<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { page } from '$app/stores';
	import { cn } from '$lib/utils';
	import { isScrolled } from './navbar-scroll.store';

	interface Tab {
		name: string;
		href: string;
	}

	export let tabs: Tab[];

	let tabsRef: (HTMLAnchorElement | null)[] = [];
	let tabsContainerRef: HTMLDivElement;
	let indicatorStyle = { left: 0, width: 0 };
	let direction: 'left' | 'right' = 'right';

	$: pathname = $page.url.pathname;
	$: activeTab = tabs.find((tab) => tab.href === pathname)?.name || tabs[0]?.name;

	async function updateIndicator() {
		await tick();
		const activeIndex = tabs.findIndex((tab) => tab.name === activeTab);
		const activeElement = tabsRef[activeIndex];

		if (activeElement && tabsContainerRef) {
			const containerRect = tabsContainerRef.getBoundingClientRect();
			const elementRect = activeElement.getBoundingClientRect();

			indicatorStyle = {
				left: elementRect.left - containerRect.left + tabsContainerRef.scrollLeft,
				width: activeElement.offsetWidth
			};
		}
	}

	function handleTabClick(tabName: string) {
		const currentIndex = tabs.findIndex((tab) => tab.name === activeTab);
		const newIndex = tabs.findIndex((tab) => tab.name === tabName);
		direction = newIndex > currentIndex ? 'right' : 'left';
	}

	$: if (activeTab || $isScrolled) {
		updateIndicator();
	}

	onMount(() => {
		updateIndicator();
	});
</script>

<div
	class={cn(
		'sticky left-0 right-0 z-40 bg-secondary-new border-b border-secondary-foreground-new transition-all duration-300',
		$isScrolled ? 'top-0' : 'top-14'
	)}
>
	<div class="mx-auto md:px-6">
		<div
			class={cn(
				'flex items-center justify-between transition-all duration-300',
				$isScrolled ? 'h-14' : 'h-12'
			)}
		>
			<div class="flex items-center gap-4 flex-1 min-w-0 h-full">
				<div
					class={cn(
						'hidden md:flex items-center gap-4 transition-all duration-300 ease-in-out flex-shrink-0',
						$isScrolled ? 'w-[44px]' : 'w-0'
					)}
				></div>
				<div
					bind:this={tabsContainerRef}
					class="relative flex-1 overflow-x-auto scrollbar-hide min-w-0 h-full"
					style="scrollbar-width: none; -ms-overflow-style: none;"
				>
					<div class="relative flex items-center gap-1 w-max h-full">
						{#each tabs as tab, index}
							<a
								href={tab.href}
								bind:this={tabsRef[index]}
								on:click={() => handleTabClick(tab.name)}
								class={cn(
									'relative px-3 py-2 text-sm font-medium transition-colors rounded-md whitespace-nowrap',
									activeTab === tab.name
										? 'text-foreground'
										: 'text-muted-foreground hover:text-foreground hover:bg-muted/50'
								)}
							>
								{tab.name}
							</a>
						{/each}
						<div
							class="absolute bottom-[-1px] h-[3px] bg-foreground transition-all duration-300 ease-out"
							style="left: {indicatorStyle.left}px; width: {indicatorStyle.width}px; transform-origin: {direction ===
							'right'
								? 'left'
								: 'right'};"
						></div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	.scrollbar-hide::-webkit-scrollbar {
		display: none;
	}
</style>
