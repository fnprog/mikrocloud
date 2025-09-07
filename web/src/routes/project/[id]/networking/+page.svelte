<script>
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import {
		Search,
		HelpCircle,
		ChevronDown,
		ArrowLeft,
		Plus,
		Network,
		Shield,
		Globe,
		Lock,
		Unlock,
		Edit,
		Trash2,
		Save,
		X,
		AlertTriangle
	} from 'lucide-svelte';

	const projectId = page.params.id;

	let project = $state({
		name: 'Focalpoint Dashboard',
		workspace: 'Focalpoint',
		category: 'Applications'
	});

	let showAddRuleModal = $state(false);
	let editingRule = $state(null);

	let newRule = $state({
		name: '',
		type: 'allow',
		protocol: 'tcp',
		port: '',
		source: '',
		description: ''
	});

	// Mock networking data
	let networkConfig = $state({
		publicIP: '203.0.113.42',
		privateIP: '10.0.1.15',
		region: 'us-east-1',
		vpc: 'vpc-12345678',
		subnet: 'subnet-87654321'
	});

	let firewallRules = $state([
		{
			id: 1,
			name: 'HTTP Traffic',
			type: 'allow',
			protocol: 'tcp',
			port: '80',
			source: '0.0.0.0/0',
			description: 'Allow HTTP traffic from anywhere',
			enabled: true,
			createdAt: '2024-01-15'
		},
		{
			id: 2,
			name: 'HTTPS Traffic',
			type: 'allow',
			protocol: 'tcp',
			port: '443',
			source: '0.0.0.0/0',
			description: 'Allow HTTPS traffic from anywhere',
			enabled: true,
			createdAt: '2024-01-15'
		},
		{
			id: 3,
			name: 'SSH Access',
			type: 'allow',
			protocol: 'tcp',
			port: '22',
			source: '203.0.113.0/24',
			description: 'SSH access from office network',
			enabled: true,
			createdAt: '2024-01-20'
		},
		{
			id: 4,
			name: 'Block Suspicious IPs',
			type: 'deny',
			protocol: 'all',
			port: 'all',
			source: '192.0.2.0/24',
			description: 'Block known malicious IP range',
			enabled: true,
			createdAt: '2024-01-25'
		}
	]);

	let loadBalancer = $state({
		enabled: true,
		algorithm: 'round_robin',
		healthCheck: {
			path: '/health',
			interval: 30,
			timeout: 5,
			unhealthyThreshold: 3
		},
		targets: [
			{ id: 1, ip: '10.0.1.10', port: 3000, status: 'healthy', weight: 100 },
			{ id: 2, ip: '10.0.1.11', port: 3000, status: 'healthy', weight: 100 },
			{ id: 3, ip: '10.0.1.12', port: 3000, status: 'unhealthy', weight: 0 }
		]
	});

	function getRuleTypeColor(type) {
		return type === 'allow' ? 'text-green-600 bg-green-50' : 'text-red-600 bg-red-50';
	}

	function getTargetStatusColor(status) {
		switch (status) {
			case 'healthy':
				return 'text-green-500';
			case 'unhealthy':
				return 'text-red-500';
			case 'draining':
				return 'text-yellow-500';
			default:
				return 'text-gray-500';
		}
	}

	function addFirewallRule() {
		if (newRule.name && newRule.port && newRule.source) {
			const newId = Math.max(...firewallRules.map((r) => r.id)) + 1;
			firewallRules = [
				...firewallRules,
				{
					id: newId,
					...newRule,
					enabled: true,
					createdAt: new Date().toISOString().split('T')[0]
				}
			];
			newRule = { name: '', type: 'allow', protocol: 'tcp', port: '', source: '', description: '' };
			showAddRuleModal = false;
		}
	}

	function toggleRule(ruleId) {
		const rule = firewallRules.find((r) => r.id === ruleId);
		if (rule) {
			rule.enabled = !rule.enabled;
		}
	}

	function deleteRule(ruleId) {
		firewallRules = firewallRules.filter((r) => r.id !== ruleId);
	}

	function startEditing(rule) {
		editingRule = { ...rule };
	}

	function saveEdit() {
		const index = firewallRules.findIndex((r) => r.id === editingRule.id);
		if (index !== -1) {
			firewallRules[index] = editingRule;
		}
		editingRule = null;
	}

	function cancelEditing() {
		editingRule = null;
	}
</script>

<svelte:head>
	<title>Networking - {project.name}</title>
</svelte:head>

<div class="flex-1 p-6">
	<!-- Header -->
	<div class="flex items-center justify-between mb-6">
		<div class="flex items-center space-x-4">
			<div>
				<h1 class="text-2xl font-semibold text-gray-900">Networking</h1>
				<p class="text-sm text-gray-500 mt-1">
					Configure network settings, firewall rules, and load balancing for {project.name}.
				</p>
			</div>
		</div>
	</div>

	<!-- Network Overview -->
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
		<Card>
			<CardContent class="p-6">
				<div class="flex items-center space-x-3">
					<Globe class="w-8 h-8 text-blue-500" />
					<div>
						<p class="text-sm font-medium text-gray-600">Public IP</p>
						<p class="text-lg font-bold text-gray-900">{networkConfig.publicIP}</p>
					</div>
				</div>
			</CardContent>
		</Card>

		<Card>
			<CardContent class="p-6">
				<div class="flex items-center space-x-3">
					<Network class="w-8 h-8 text-green-500" />
					<div>
						<p class="text-sm font-medium text-gray-600">Private IP</p>
						<p class="text-lg font-bold text-gray-900">{networkConfig.privateIP}</p>
					</div>
				</div>
			</CardContent>
		</Card>

		<Card>
			<CardContent class="p-6">
				<div class="flex items-center space-x-3">
					<Shield class="w-8 h-8 text-purple-500" />
					<div>
						<p class="text-sm font-medium text-gray-600">VPC</p>
						<p class="text-lg font-bold text-gray-900">{networkConfig.vpc}</p>
					</div>
				</div>
			</CardContent>
		</Card>

		<Card>
			<CardContent class="p-6">
				<div class="flex items-center space-x-3">
					<Network class="w-8 h-8 text-orange-500" />
					<div>
						<p class="text-sm font-medium text-gray-600">Subnet</p>
						<p class="text-lg font-bold text-gray-900">{networkConfig.subnet}</p>
					</div>
				</div>
			</CardContent>
		</Card>
	</div>

	<!-- Firewall Rules -->
	<Card class="mb-8">
		<CardHeader>
			<div class="flex items-center justify-between">
				<CardTitle>Firewall Rules</CardTitle>
				<Button onclick={() => (showAddRuleModal = true)}>
					<Plus class="w-4 h-4 mr-2" />
					Add Rule
				</Button>
			</div>
		</CardHeader>
		<CardContent>
			<div class="space-y-4">
				{#each firewallRules as rule (rule.id)}
					<div class="border rounded-lg p-4">
						{#if editingRule && editingRule.id === rule.id}
							<!-- Edit Mode -->
							<div class="grid grid-cols-2 gap-4">
								<div>
									<Label for="edit-name">Rule Name</Label>
									<Input id="edit-name" bind:value={editingRule.name} />
								</div>
								<div>
									<Label for="edit-type">Type</Label>
									<select
										id="edit-type"
										bind:value={editingRule.type}
										class="w-full px-3 py-2 border border-gray-300 rounded-md"
									>
										<option value="allow">Allow</option>
										<option value="deny">Deny</option>
									</select>
								</div>
								<div>
									<Label for="edit-protocol">Protocol</Label>
									<select
										id="edit-protocol"
										bind:value={editingRule.protocol}
										class="w-full px-3 py-2 border border-gray-300 rounded-md"
									>
										<option value="tcp">TCP</option>
										<option value="udp">UDP</option>
										<option value="all">All</option>
									</select>
								</div>
								<div>
									<Label for="edit-port">Port</Label>
									<Input id="edit-port" bind:value={editingRule.port} placeholder="80, 443, 22" />
								</div>
								<div class="col-span-2">
									<Label for="edit-source">Source</Label>
									<Input id="edit-source" bind:value={editingRule.source} placeholder="0.0.0.0/0" />
								</div>
								<div class="col-span-2">
									<Label for="edit-description">Description</Label>
									<Input id="edit-description" bind:value={editingRule.description} />
								</div>
								<div class="col-span-2 flex space-x-2">
									<Button size="sm" onclick={saveEdit}>
										<Save class="w-4 h-4 mr-1" />
										Save
									</Button>
									<Button size="sm" variant="outline" onclick={cancelEditing}>
										<X class="w-4 h-4 mr-1" />
										Cancel
									</Button>
								</div>
							</div>
						{:else}
							<!-- View Mode -->
							<div class="flex items-center justify-between">
								<div class="flex-1">
									<div class="flex items-center space-x-3 mb-2">
										<h3 class="font-medium text-gray-900">{rule.name}</h3>
										<Badge variant="outline" class="text-xs {getRuleTypeColor(rule.type)}">
											{rule.type.toUpperCase()}
										</Badge>
										<button onclick={() => toggleRule(rule.id)} class="flex items-center space-x-1">
											{#if rule.enabled}
												<Unlock class="w-4 h-4 text-green-500" />
												<span class="text-sm text-green-600">Enabled</span>
											{:else}
												<Lock class="w-4 h-4 text-red-500" />
												<span class="text-sm text-red-600">Disabled</span>
											{/if}
										</button>
									</div>
									<div class="grid grid-cols-4 gap-4 text-sm text-gray-600">
										<div>
											<span class="font-medium">Protocol:</span>
											{rule.protocol.toUpperCase()}
										</div>
										<div>
											<span class="font-medium">Port:</span>
											{rule.port}
										</div>
										<div>
											<span class="font-medium">Source:</span>
											{rule.source}
										</div>
										<div>
											<span class="font-medium">Created:</span>
											{rule.createdAt}
										</div>
									</div>
									{#if rule.description}
										<p class="text-sm text-gray-600 mt-2">{rule.description}</p>
									{/if}
								</div>
								<div class="flex items-center space-x-2">
									<Button size="sm" variant="outline" onclick={() => startEditing(rule)}>
										<Edit class="w-4 h-4" />
									</Button>
									<Button size="sm" variant="outline" onclick={() => deleteRule(rule.id)}>
										<Trash2 class="w-4 h-4" />
									</Button>
								</div>
							</div>
						{/if}
					</div>
				{/each}
			</div>
		</CardContent>
	</Card>

	<!-- Load Balancer -->
	<Card>
		<CardHeader>
			<CardTitle>Load Balancer</CardTitle>
		</CardHeader>
		<CardContent>
			<div class="space-y-6">
				<!-- Load Balancer Settings -->
				<div class="grid grid-cols-2 gap-6">
					<div>
						<Label>Algorithm</Label>
						<select
							bind:value={loadBalancer.algorithm}
							class="w-full px-3 py-2 border border-gray-300 rounded-md mt-1"
						>
							<option value="round_robin">Round Robin</option>
							<option value="least_connections">Least Connections</option>
							<option value="ip_hash">IP Hash</option>
						</select>
					</div>
					<div>
						<Label>Health Check Path</Label>
						<Input bind:value={loadBalancer.healthCheck.path} class="mt-1" />
					</div>
				</div>

				<!-- Targets -->
				<div>
					<h4 class="font-medium text-gray-900 mb-4">Targets</h4>
					<div class="space-y-3">
						{#each loadBalancer.targets as target (target.id)}
							<div class="flex items-center justify-between p-3 border rounded-lg">
								<div class="flex items-center space-x-4">
									<div
										class="w-3 h-3 rounded-full {target.status === 'healthy'
											? 'bg-green-500'
											: target.status === 'unhealthy'
												? 'bg-red-500'
												: 'bg-yellow-500'}"
									></div>
									<div>
										<p class="font-medium text-gray-900">{target.ip}:{target.port}</p>
										<p class="text-sm {getTargetStatusColor(target.status)}">{target.status}</p>
									</div>
								</div>
								<div class="flex items-center space-x-4">
									<div class="text-sm text-gray-600">
										Weight: {target.weight}%
									</div>
									<Button size="sm" variant="outline">
										<Edit class="w-4 h-4" />
									</Button>
								</div>
							</div>
						{/each}
					</div>
				</div>
			</div>
		</CardContent>
	</Card>
</div>

<!-- Add Firewall Rule Modal -->
{#if showAddRuleModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<Card class="w-full max-w-lg">
			<CardHeader>
				<CardTitle>Add Firewall Rule</CardTitle>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="grid grid-cols-2 gap-4">
					<div>
						<Label for="name">Rule Name</Label>
						<Input id="name" bind:value={newRule.name} placeholder="HTTP Traffic" />
					</div>
					<div>
						<Label for="type">Type</Label>
						<select
							id="type"
							bind:value={newRule.type}
							class="w-full px-3 py-2 border border-gray-300 rounded-md"
						>
							<option value="allow">Allow</option>
							<option value="deny">Deny</option>
						</select>
					</div>
					<div>
						<Label for="protocol">Protocol</Label>
						<select
							id="protocol"
							bind:value={newRule.protocol}
							class="w-full px-3 py-2 border border-gray-300 rounded-md"
						>
							<option value="tcp">TCP</option>
							<option value="udp">UDP</option>
							<option value="all">All</option>
						</select>
					</div>
					<div>
						<Label for="port">Port</Label>
						<Input id="port" bind:value={newRule.port} placeholder="80, 443, 22" />
					</div>
				</div>
				<div>
					<Label for="source">Source</Label>
					<Input id="source" bind:value={newRule.source} placeholder="0.0.0.0/0" />
				</div>
				<div>
					<Label for="description">Description</Label>
					<Input
						id="description"
						bind:value={newRule.description}
						placeholder="Allow HTTP traffic from anywhere"
					/>
				</div>
				<div class="flex space-x-2 pt-4">
					<Button onclick={addFirewallRule} class="flex-1">Add Rule</Button>
					<Button variant="outline" onclick={() => (showAddRuleModal = false)} class="flex-1"
						>Cancel</Button
					>
				</div>
			</CardContent>
		</Card>
	</div>
{/if}
