<script>
	import { page } from '$app/stores';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Database, Plus, Table2, Search, Play, Settings } from 'lucide-svelte';

	let projectId = $derived($page.params.id);
	let envId = $derived($page.params.env_id);
	let resId = $derived($page.params.res_id);
</script>

<svelte:head>
	<title>Database Studio</title>
</svelte:head>

<div class="flex-1 p-6">
	<div class="flex items-center justify-between mb-6">
		<div class="flex items-center space-x-4">
			<div>
				<h1 class="text-2xl font-semibold text-gray-900">Database Studio</h1>
				<p class="text-sm text-gray-500 mt-1">
					Visual database management interface for querying and managing your data.
				</p>
			</div>
		</div>
	</div>

	<Card>
		<CardHeader>
			<div class="flex items-center justify-between">
				<div class="flex items-center space-x-2">
					<Database class="w-5 h-5 text-gray-600" />
					<CardTitle>Database Studio - Coming Soon</CardTitle>
				</div>
			</div>
		</CardHeader>
		<CardContent>
			<div class="text-center py-12 space-y-4">
				<div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-blue-100">
					<Database class="w-8 h-8 text-blue-600" />
				</div>
				<div>
					<h3 class="text-lg font-semibold text-gray-900 mb-2">
						Advanced Database Management Interface
					</h3>
					<p class="text-sm text-gray-600 max-w-2xl mx-auto mb-4">
						The Database Studio will provide a comprehensive visual interface for managing your
						database, similar to tools like Supabase Table Editor, phpMyAdmin, and MongoDB Compass.
					</p>
				</div>
				<div class="grid grid-cols-2 md:grid-cols-4 gap-4 max-w-4xl mx-auto mt-8">
					<div class="p-4 bg-gray-50 rounded-lg">
						<Table2 class="w-6 h-6 text-gray-700 mb-2 mx-auto" />
						<h4 class="font-medium text-sm text-gray-900 mb-1">Table Browser</h4>
						<p class="text-xs text-gray-600">
							Browse, filter, and paginate through database tables
						</p>
					</div>
					<div class="p-4 bg-gray-50 rounded-lg">
						<Search class="w-6 h-6 text-gray-700 mb-2 mx-auto" />
						<h4 class="font-medium text-sm text-gray-900 mb-1">Query Builder</h4>
						<p class="text-xs text-gray-600">Visual query builder with SQL preview</p>
					</div>
					<div class="p-4 bg-gray-50 rounded-lg">
						<Play class="w-6 h-6 text-gray-700 mb-2 mx-auto" />
						<h4 class="font-medium text-sm text-gray-900 mb-1">SQL Console</h4>
						<p class="text-xs text-gray-600">
							Execute custom SQL queries with syntax highlighting
						</p>
					</div>
					<div class="p-4 bg-gray-50 rounded-lg">
						<Plus class="w-6 h-6 text-gray-700 mb-2 mx-auto" />
						<h4 class="font-medium text-sm text-gray-900 mb-1">CRUD Operations</h4>
						<p class="text-xs text-gray-600">Create, read, update, and delete records easily</p>
					</div>
				</div>
				<div class="mt-8 pt-8 border-t border-gray-200">
					<h4 class="font-medium text-sm text-gray-900 mb-3">Supported Database Systems</h4>
					<div class="flex flex-wrap justify-center gap-3">
						<span class="px-3 py-1 bg-blue-100 text-blue-700 text-xs font-medium rounded-full"
							>PostgreSQL</span
						>
						<span class="px-3 py-1 bg-orange-100 text-orange-700 text-xs font-medium rounded-full"
							>MySQL</span
						>
						<span class="px-3 py-1 bg-teal-100 text-teal-700 text-xs font-medium rounded-full"
							>MariaDB</span
						>
						<span class="px-3 py-1 bg-green-100 text-green-700 text-xs font-medium rounded-full"
							>MongoDB</span
						>
						<span class="px-3 py-1 bg-red-100 text-red-700 text-xs font-medium rounded-full"
							>Redis</span
						>
						<span class="px-3 py-1 bg-yellow-100 text-yellow-700 text-xs font-medium rounded-full"
							>ClickHouse</span
						>
					</div>
				</div>
			</div>
		</CardContent>
	</Card>
</div>
