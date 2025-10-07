<script lang="ts">
	import { cn } from '$lib/utils';
	import { goto } from '$app/navigation';
	import { authApi } from '$lib/api/auth';

	import {
		DropdownMenu,
		DropdownMenuContent,
		DropdownMenuItem,
		DropdownMenuSeparator,
		DropdownMenuTrigger
	} from '$lib/components/ui/dropdown-menu/index';
	import { Avatar, AvatarFallback } from '$lib/components/ui/avatar';
	import { Button } from '$lib/components/ui/button';

	import { LogOut, Monitor, Sun, Moon, Plus } from 'lucide-svelte';

	let theme: 'light' | 'dark' | 'system' = 'system';

	function cycleTheme() {
		const themes: Array<'light' | 'dark' | 'system'> = ['light', 'dark', 'system'];
		const currentIndex = themes.indexOf(theme);
		const nextIndex = (currentIndex + 1) % themes.length;
		theme = themes[nextIndex];
	}

	function handleLogout() {
		authApi.logout();
		goto('/login');
	}
</script>

<DropdownMenu>
	<DropdownMenuTrigger>
		<Button variant="ghost" size="icon" class="rounded-full h-8 w-8">
			<Avatar class="h-8 w-8">
				<AvatarFallback class="bg-orange-500 text-white text-xs">GN</AvatarFallback>
			</Avatar>
		</Button>
	</DropdownMenuTrigger>

	<DropdownMenuContent align="end" class="w-[280px] p-2">
		<!-- Header -->
		<div class="px-2 py-3 border-b border-border mb-2">
			<div class="font-medium text-sm">gamernewone</div>
			<div class="text-xs text-muted-foreground">rickrichard80@gmail.com</div>
		</div>

		<!-- Dashboard -->
		<DropdownMenuItem class="cursor-pointer">
			<span class="text-sm">Dashboard</span>
		</DropdownMenuItem>

		<!-- Account Settings -->
		<DropdownMenuItem class="cursor-pointer">
			<span class="text-sm">Account Settings</span>
		</DropdownMenuItem>

		<!-- Create Team -->
		<DropdownMenuItem class="cursor-pointer justify-between">
			<span class="text-sm">Create Team</span>
			<Plus class="h-4 w-4 text-muted-foreground" />
		</DropdownMenuItem>

		<DropdownMenuSeparator />

		<!-- Command Menu -->
		<DropdownMenuItem class="cursor-pointer justify-between">
			<div class="flex items-center">
				<span class="text-sm">Command Menu</span>
			</div>
			<div class="flex items-center gap-1 text-xs text-muted-foreground">
				<kbd class="px-1.5 py-0.5 bg-muted rounded text-[10px]">Ctl</kbd>
				<kbd class="px-1.5 py-0.5 bg-muted rounded text-[10px]">K</kbd>
			</div>
		</DropdownMenuItem>

		<!-- Theme Switch -->
		<DropdownMenuItem class="cursor-pointer justify-between" onSelect={cycleTheme}>
			<span class="text-sm">Theme</span>
			<div class="flex items-center gap-1">
				<Monitor
					class={cn('h-4 w-4', theme === 'system' ? 'text-foreground' : 'text-muted-foreground')}
				/>
				<Sun
					class={cn('h-4 w-4', theme === 'light' ? 'text-foreground' : 'text-muted-foreground')}
				/>
				<Moon
					class={cn('h-4 w-4', theme === 'dark' ? 'text-foreground' : 'text-muted-foreground')}
				/>
			</div>
		</DropdownMenuItem>

		<DropdownMenuSeparator />

		<!-- Logout -->
		<DropdownMenuItem class="cursor-pointer justify-between" onSelect={handleLogout}>
			<span class="text-sm">Log Out</span>
			<LogOut class="h-4 w-4" />
		</DropdownMenuItem>
	</DropdownMenuContent>
</DropdownMenu>
