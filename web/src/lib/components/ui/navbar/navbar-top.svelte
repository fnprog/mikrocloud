<script lang="ts">
	import Logo from '$lib/components/logo/logo.svelte';
	import * as Avatar from '$lib/components/ui/avatar/index';
	import { isScrolled } from './navbar-scroll.store';
	import { cn } from '$lib/utils';
	import UserMenu from './nav-top-menu.svelte';

	let switcherOpen = false;
</script>

<header class="fixed top-0 left-0 right-0 z-50 bg-secondary-new">
	<div class="px-6 relative">
		<Logo
			class={cn(
				'h-[25px] absolute left-6 hidden md:flex items-center gap-4 transition-all duration-300 ease-in-out z-10',
				$isScrolled ? 'translate-y-[21px] h-[20px]' : 'translate-y-[14px]'
			)}
		/>

		<div
			class={cn(
				'flex items-center gap-4 transition-all duration-300 ease-in-out overflow-hidden',
				$isScrolled ? 'h-0 opacity-0' : 'h-14 opacity-100'
			)}
		>
			<div class="w-[25px] flex-shrink-0 hidden md:block"></div>

			<button
				on:click={() => (switcherOpen = true)}
				class="flex items-center gap-2 text-sm hover:bg-muted px-2 py-1 rounded-md transition-colors"
			>
				<span class="text-muted-foreground">/</span>
				<div class="flex items-center gap-2">
					<Avatar.Root>
						<Avatar.Image src="https://github.com/shadcn.png" alt="@shadcn" />
						<Avatar.Fallback>CN</Avatar.Fallback>
					</Avatar.Root>
					<span class="font-medium hidden md:inline">gxuri's projects</span>
				</div>
				<span class="text-muted-foreground hidden md:inline">Hobby</span>
				<svg width="12" height="12" viewBox="0 0 12 12" fill="none" class="text-muted-foreground">
					<path
						d="M3 4.5L6 7.5L9 4.5"
						stroke="currentColor"
						stroke-width="1.5"
						stroke-linecap="round"
						stroke-linejoin="round"
					/>
				</svg>
			</button>

			<!-- <button -->
			<!-- 	class="flex items-center gap-2 text-sm hover:bg-muted px-2 py-1 rounded-md transition-colors" -->
			<!-- > -->
			<!-- 	<div class="flex items-center gap-2"> -->
			<!-- 		<Avatar.Root> -->
			<!-- 			<Avatar.Image src="https://github.com/shadcn.png" alt="@shadcn" /> -->
			<!-- 			<Avatar.Fallback>CN</Avatar.Fallback> -->
			<!-- 		</Avatar.Root> -->
			<!-- 		<span class="font-medium">skiper-ui-demo</span> -->
			<!-- 	</div> -->
			<!-- 	<svg width="12" height="12" viewBox="0 0 12 12" fill="none" class="text-muted-foreground"> -->
			<!-- 		<path -->
			<!-- 			d="M3 4.5L6 7.5L9 4.5" -->
			<!-- 			stroke="currentColor" -->
			<!-- 			stroke-width="1.5" -->
			<!-- 			stroke-linecap="round" -->
			<!-- 			stroke-linejoin="round" -->
			<!-- 		/> -->
			<!-- 	</svg> -->
			<!-- </button> -->
			<div class="ml-auto">
				<UserMenu />
			</div>
		</div>
	</div>

	<!-- <Avatar.Root> -->
	<!-- 	<Avatar.Image src="https://github.com/shadcn.png" alt="@shadcn" /> -->
	<!-- 	<Avatar.Fallback>CN</Avatar.Fallback> -->
	<!-- </Avatar.Root> -->
</header>
