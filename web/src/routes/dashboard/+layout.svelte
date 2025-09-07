<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import Grid_3x3 from 'lucide-svelte/icons/grid-3x3';
	import Logo from '$lib/components/logo/Logo.svelte';
	import Grid_2x2Check from 'lucide-svelte/icons/grid-2x2-check';
	import {
		Bell,
		FileText,
		GitBranch,
		HardDrive,
		Key,
		LogOut,
		Server,
		Terminal,
		User
	} from 'lucide-svelte';

	let { children } = $props();

	const items = [
		{
			item: 'Dashboard',
			route: '/dashboard',
			icon: Grid_3x3
		},
		{
			item: 'Services',
			route: '/dashboard/services',
			icon: Grid_2x2Check
		},
		{
			item: 'Servers',
			route: '/dashboard/servers',
			icon: Server
		},
		{
			item: 'Sources',
			route: '/dashboard/servers',
			icon: GitBranch
		},
		{
			item: 'Backups',
			route: '/dashboard/servers',
			icon: HardDrive
		},
		{
			item: 'Shared Variables',
			route: '/dashboard/env',
			icon: FileText
		},
		{
			item: 'Notifications',
			route: '/dashboard/env',
			icon: Bell
		},
		{
			item: 'Keys and Token',
			route: '/dashboard/keys',
			icon: Key
		},
		{
			item: 'Terminal',
			route: '/dashboard/term',
			icon: Terminal
		}
	];
</script>

<Sidebar.Provider>
	<Sidebar.Root collapsible="icon" class="bg-[#EEEEEE]">
		<Sidebar.Header>
			<div class="p-4 border-b border-gray-200">
				<div class="flex items-center space-x-3">
					<Logo class="w-8 h-8" />
				</div>
			</div>
		</Sidebar.Header>
		<Sidebar.Content class="p-4">
			<Sidebar.Menu>
				{#each items as item (item.item)}
					<Sidebar.MenuItem>
						<Sidebar.MenuButton tooltipContent={item.item}>
							{#snippet child({ props })}
								<a href={item.route} {...props}>
									<item.icon class="w-5 h-5 mr-3" />
									<span>{item.item}</span>
								</a>
							{/snippet}
						</Sidebar.MenuButton>
					</Sidebar.MenuItem>
				{/each}
			</Sidebar.Menu>
		</Sidebar.Content>
		<Sidebar.Footer>
			<Sidebar.Menu>
				<Sidebar.MenuItem>
					<Sidebar.MenuButton tooltipContent="Profile">
						<a href="/dashboard/profile">
							<User class="w-5 h-5 mr-3" />
							<span>Profile</span>
						</a>
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
				<Sidebar.MenuItem>
					<Sidebar.MenuButton tooltipContent="Team">
						<a href="/dashboard/team">
							<User class="w-5 h-5 mr-3" />
							<span>Team</span>
						</a>
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
				<Sidebar.MenuItem>
					<Sidebar.MenuButton tooltipContent="Settings">
						<a href="/dashboard/settings">
							<User class="w-5 h-5 mr-3" />
							<span>Settings</span>
						</a>
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
				<Sidebar.MenuItem>
					<Sidebar.MenuButton
						class="w-full flex items-center px-3 py-2 text-sm font-medium text-gray-600 rounded-md hover:bg-gray-50 hover:text-gray-900"
						tooltipContent="Logout"
					>
						<LogOut class="w-5 h-5 mr-3" />
						Logout
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			</Sidebar.Menu>
		</Sidebar.Footer>
		<Sidebar.Rail />
	</Sidebar.Root>
	<Sidebar.Inset class="bg-[#F3F3F3]">
		{@render children()}
	</Sidebar.Inset>
</Sidebar.Provider>
