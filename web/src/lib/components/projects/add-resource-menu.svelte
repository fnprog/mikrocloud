<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import {
		DropdownMenu,
		DropdownMenuContent,
		DropdownMenuItem,
		DropdownMenuTrigger
	} from '$lib/components/ui/dropdown-menu';
	import { Plus, Server, Database, LayoutGrid } from 'lucide-svelte';

	interface Props {
		projectId: string;
		envId?: string;
	}

	let { projectId, envId }: Props = $props();
</script>

<DropdownMenu>
	<DropdownMenuTrigger>
		<Button variant="default">
			<Plus class="h-4 w-4" />
			Add Resource
		</Button>
	</DropdownMenuTrigger>
	<DropdownMenuContent align="end" class="w-56">
		<a href={`/dashboard/project/${projectId}/${envId}/create-app`}>
			<DropdownMenuItem>
				<Server class="mr-2 h-4 w-4" />
				Application
			</DropdownMenuItem>
		</a>
		<a href={`/dashboard/project/${projectId}/${envId}/create-db`}>
			<DropdownMenuItem>
				<Database class="mr-2 h-4 w-4" />
				Database
			</DropdownMenuItem>
		</a>
		<a href={`/dashboard/project/${projectId}/${envId}/create-service`}>
			<DropdownMenuItem>
				<LayoutGrid class="mr-2 h-4 w-4" />
				Template
			</DropdownMenuItem>
		</a>
	</DropdownMenuContent>
</DropdownMenu>
