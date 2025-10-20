<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import { Progress } from '$lib/components/ui/progress';
	import {
		Activity,
		AlertTriangle,
		ChevronDown,
		ChevronUp,
		Clock,
		Database,
		Download,
		RefreshCw,
		Search,
		TrendingUp,
		Zap
	} from 'lucide-svelte';
	import AreaChart from '$lib/components/charts/area-chart.svelte';
	import BarChart from '$lib/components/charts/bar-chart.svelte';
	import LineChart from '$lib/components/charts/line-chart.svelte';

	interface Trace {
		id: string;
		endpoint: string;
		method: string;
		statusCode: number;
		duration: number;
		timestamp: Date;
		errorMessage?: string;
		spans: TraceSpan[];
	}

	interface TraceSpan {
		id: string;
		name: string;
		type: 'http' | 'database' | 'cache' | 'external' | 'internal';
		duration: number;
		startOffset: number;
		status: 'success' | 'error' | 'warning';
		details?: string;
	}

	interface EndpointMetrics {
		endpoint: string;
		method: string;
		requestCount: number;
		errorCount: number;
		errorRate: number;
		avgDuration: number;
		p50: number;
		p95: number;
		p99: number;
		slowestTrace?: string;
	}

	interface SlowQuery {
		id: string;
		query: string;
		database: string;
		duration: number;
		timestamp: Date;
		executionCount: number;
		affectedRows: number;
		traceId: string;
	}

	const mockTraces: Trace[] = [
		{
			id: 'trace-1',
			endpoint: '/api/deployments',
			method: 'POST',
			statusCode: 201,
			duration: 847,
			timestamp: new Date(Date.now() - 120000),
			spans: [
				{
					id: 'span-1',
					name: 'HTTP Request',
					type: 'http',
					duration: 847,
					startOffset: 0,
					status: 'success'
				},
				{
					id: 'span-2',
					name: 'Auth Validation',
					type: 'internal',
					duration: 45,
					startOffset: 12,
					status: 'success'
				},
				{
					id: 'span-3',
					name: 'Database Query',
					type: 'database',
					duration: 234,
					startOffset: 78,
					status: 'success',
					details: 'SELECT * FROM deployments WHERE id = $1'
				},
				{
					id: 'span-4',
					name: 'Container Build',
					type: 'external',
					duration: 456,
					startOffset: 340,
					status: 'success',
					details: 'Docker Build API'
				},
				{
					id: 'span-5',
					name: 'Cache Write',
					type: 'cache',
					duration: 23,
					startOffset: 812,
					status: 'success'
				}
			]
		},
		{
			id: 'trace-2',
			endpoint: '/api/applications/:id',
			method: 'GET',
			statusCode: 200,
			duration: 156,
			timestamp: new Date(Date.now() - 45000),
			spans: [
				{
					id: 'span-6',
					name: 'HTTP Request',
					type: 'http',
					duration: 156,
					startOffset: 0,
					status: 'success'
				},
				{
					id: 'span-7',
					name: 'Cache Read',
					type: 'cache',
					duration: 12,
					startOffset: 8,
					status: 'success'
				},
				{
					id: 'span-8',
					name: 'Database Query',
					type: 'database',
					duration: 89,
					startOffset: 34,
					status: 'success',
					details: 'SELECT * FROM applications WHERE id = $1'
				}
			]
		},
		{
			id: 'trace-3',
			endpoint: '/api/databases/connect',
			method: 'POST',
			statusCode: 500,
			duration: 5234,
			timestamp: new Date(Date.now() - 300000),
			errorMessage: 'Connection timeout to database server',
			spans: [
				{
					id: 'span-9',
					name: 'HTTP Request',
					type: 'http',
					duration: 5234,
					startOffset: 0,
					status: 'error'
				},
				{
					id: 'span-10',
					name: 'Database Connection',
					type: 'database',
					duration: 5000,
					startOffset: 123,
					status: 'error',
					details: 'Timeout connecting to postgres://db.example.com'
				}
			]
		},
		{
			id: 'trace-4',
			endpoint: '/api/logs/stream',
			method: 'GET',
			statusCode: 200,
			duration: 234,
			timestamp: new Date(Date.now() - 90000),
			spans: [
				{
					id: 'span-11',
					name: 'HTTP Request',
					type: 'http',
					duration: 234,
					startOffset: 0,
					status: 'success'
				},
				{
					id: 'span-12',
					name: 'Database Query',
					type: 'database',
					duration: 167,
					startOffset: 45,
					status: 'warning',
					details: 'SELECT * FROM logs ORDER BY timestamp DESC LIMIT 1000'
				}
			]
		},
		{
			id: 'trace-5',
			endpoint: '/api/users/:id',
			method: 'PUT',
			statusCode: 200,
			duration: 312,
			timestamp: new Date(Date.now() - 180000),
			spans: [
				{
					id: 'span-13',
					name: 'HTTP Request',
					type: 'http',
					duration: 312,
					startOffset: 0,
					status: 'success'
				},
				{
					id: 'span-14',
					name: 'Auth Check',
					type: 'internal',
					duration: 34,
					startOffset: 12,
					status: 'success'
				},
				{
					id: 'span-15',
					name: 'Database Update',
					type: 'database',
					duration: 145,
					startOffset: 78,
					status: 'success',
					details: 'UPDATE users SET name = $1 WHERE id = $2'
				},
				{
					id: 'span-16',
					name: 'Cache Invalidate',
					type: 'cache',
					duration: 18,
					startOffset: 245,
					status: 'success'
				}
			]
		}
	];

	const mockEndpointMetrics: EndpointMetrics[] = [
		{
			endpoint: '/api/deployments',
			method: 'POST',
			requestCount: 1247,
			errorCount: 12,
			errorRate: 0.96,
			avgDuration: 892,
			p50: 678,
			p95: 1234,
			p99: 2145,
			slowestTrace: 'trace-1'
		},
		{
			endpoint: '/api/applications/:id',
			method: 'GET',
			requestCount: 8934,
			errorCount: 45,
			errorRate: 0.5,
			avgDuration: 134,
			p50: 98,
			p95: 287,
			p99: 456,
			slowestTrace: 'trace-2'
		},
		{
			endpoint: '/api/databases/connect',
			method: 'POST',
			requestCount: 456,
			errorCount: 89,
			errorRate: 19.52,
			avgDuration: 2345,
			p50: 1234,
			p95: 4567,
			p99: 6789,
			slowestTrace: 'trace-3'
		},
		{
			endpoint: '/api/logs/stream',
			method: 'GET',
			requestCount: 3421,
			errorCount: 23,
			errorRate: 0.67,
			avgDuration: 245,
			p50: 189,
			p95: 456,
			p99: 678,
			slowestTrace: 'trace-4'
		},
		{
			endpoint: '/api/users/:id',
			method: 'PUT',
			requestCount: 2134,
			errorCount: 8,
			errorRate: 0.37,
			avgDuration: 298,
			p50: 234,
			p95: 512,
			p99: 789
		},
		{
			endpoint: '/api/projects',
			method: 'GET',
			requestCount: 5678,
			errorCount: 12,
			errorRate: 0.21,
			avgDuration: 167,
			p50: 123,
			p95: 289,
			p99: 445
		},
		{
			endpoint: '/api/services/:id/restart',
			method: 'POST',
			requestCount: 789,
			errorCount: 34,
			errorRate: 4.31,
			avgDuration: 1456,
			p50: 1123,
			p95: 2345,
			p99: 3456
		}
	];

	const mockSlowQueries: SlowQuery[] = [
		{
			id: 'query-1',
			query:
				'SELECT d.*, a.name, p.name FROM deployments d JOIN applications a ON d.application_id = a.id JOIN projects p ON a.project_id = p.id WHERE d.status = \'running\' ORDER BY d.created_at DESC',
			database: 'mikrocloud_main',
			duration: 3456,
			timestamp: new Date(Date.now() - 180000),
			executionCount: 234,
			affectedRows: 1247,
			traceId: 'trace-1'
		},
		{
			id: 'query-2',
			query:
				'SELECT * FROM logs WHERE application_id = $1 AND timestamp > $2 ORDER BY timestamp DESC LIMIT 1000',
			database: 'mikrocloud_analytics',
			duration: 2145,
			timestamp: new Date(Date.now() - 90000),
			executionCount: 1456,
			affectedRows: 8934,
			traceId: 'trace-4'
		},
		{
			id: 'query-3',
			query:
				'UPDATE containers SET status = $1, updated_at = NOW() WHERE id IN (SELECT id FROM containers WHERE server_id = $2)',
			database: 'mikrocloud_main',
			duration: 1789,
			timestamp: new Date(Date.now() - 240000),
			executionCount: 89,
			affectedRows: 456,
			traceId: 'trace-5'
		},
		{
			id: 'query-4',
			query:
				'SELECT COUNT(*) as total, status FROM applications GROUP BY status HAVING COUNT(*) > 10',
			database: 'mikrocloud_main',
			duration: 1234,
			timestamp: new Date(Date.now() - 360000),
			executionCount: 567,
			affectedRows: 12,
			traceId: 'trace-2'
		}
	];

	const latencyChartData = [
		{ label: '00:00', value: 145 },
		{ label: '02:00', value: 178 },
		{ label: '04:00', value: 234 },
		{ label: '06:00', value: 312 },
		{ label: '08:00', value: 456 },
		{ label: '10:00', value: 523 },
		{ label: '12:00', value: 678 },
		{ label: '14:00', value: 789 },
		{ label: '16:00', value: 612 },
		{ label: '18:00', value: 445 },
		{ label: '20:00', value: 298 },
		{ label: '22:00', value: 189 }
	];

	const requestRateChartData = [
		{ label: '00:00', value: 1234 },
		{ label: '02:00', value: 1456 },
		{ label: '04:00', value: 1678 },
		{ label: '06:00', value: 2145 },
		{ label: '08:00', value: 3456 },
		{ label: '10:00', value: 4523 },
		{ label: '12:00', value: 5678 },
		{ label: '14:00', value: 6234 },
		{ label: '16:00', value: 5123 },
		{ label: '18:00', value: 3789 },
		{ label: '20:00', value: 2456 },
		{ label: '22:00', value: 1678 }
	];

	const errorRateChartData = [
		{ label: '00:00', value: 0.5 },
		{ label: '02:00', value: 0.8 },
		{ label: '04:00', value: 1.2 },
		{ label: '06:00', value: 1.8 },
		{ label: '08:00', value: 2.3 },
		{ label: '10:00', value: 1.9 },
		{ label: '12:00', value: 2.5 },
		{ label: '14:00', value: 3.1 },
		{ label: '16:00', value: 2.2 },
		{ label: '18:00', value: 1.5 },
		{ label: '20:00', value: 0.9 },
		{ label: '22:00', value: 0.6 }
	];

	const throughputChartData = [
		{ label: 'GET', value: 12456 },
		{ label: 'POST', value: 4567 },
		{ label: 'PUT', value: 2345 },
		{ label: 'DELETE', value: 1234 },
		{ label: 'PATCH', value: 567 }
	];

	let searchQuery = $state('');
	let statusFilter = $state('all');
	let methodFilter = $state('all');
	let expandedTrace = $state<string | null>(null);
	let expandedEndpoint = $state<string | null>(null);
	let expandedQuery = $state<string | null>(null);
	let isRefreshing = $state(false);
	let activeTab = $state<'traces' | 'endpoints' | 'queries'>('traces');

	const filteredTraces = $derived(() => {
		let filtered = mockTraces;

		if (searchQuery) {
			filtered = filtered.filter(
				(trace) =>
					trace.endpoint.toLowerCase().includes(searchQuery.toLowerCase()) ||
					trace.method.toLowerCase().includes(searchQuery.toLowerCase()) ||
					trace.id.toLowerCase().includes(searchQuery.toLowerCase())
			);
		}

		if (statusFilter !== 'all') {
			if (statusFilter === 'success') {
				filtered = filtered.filter((trace) => trace.statusCode >= 200 && trace.statusCode < 300);
			} else if (statusFilter === 'error') {
				filtered = filtered.filter((trace) => trace.statusCode >= 400);
			} else if (statusFilter === 'slow') {
				filtered = filtered.filter((trace) => trace.duration > 1000);
			}
		}

		if (methodFilter !== 'all') {
			filtered = filtered.filter((trace) => trace.method === methodFilter);
		}

		return filtered;
	});

	const filteredEndpoints = $derived(() => {
		let filtered = mockEndpointMetrics;

		if (searchQuery) {
			filtered = filtered.filter(
				(endpoint) =>
					endpoint.endpoint.toLowerCase().includes(searchQuery.toLowerCase()) ||
					endpoint.method.toLowerCase().includes(searchQuery.toLowerCase())
			);
		}

		if (methodFilter !== 'all') {
			filtered = filtered.filter((endpoint) => endpoint.method === methodFilter);
		}

		if (statusFilter === 'slow') {
			filtered = filtered.filter((endpoint) => endpoint.avgDuration > 500);
		} else if (statusFilter === 'error') {
			filtered = filtered.filter((endpoint) => endpoint.errorRate > 1);
		}

		return filtered.sort((a, b) => b.p99 - a.p99);
	});

	const filteredQueries = $derived(() => {
		let filtered = mockSlowQueries;

		if (searchQuery) {
			filtered = filtered.filter(
				(query) =>
					query.query.toLowerCase().includes(searchQuery.toLowerCase()) ||
					query.database.toLowerCase().includes(searchQuery.toLowerCase()) ||
					query.id.toLowerCase().includes(searchQuery.toLowerCase())
			);
		}

		return filtered.sort((a, b) => b.duration - a.duration);
	});

	const stats = $derived(() => {
		const totalRequests = mockEndpointMetrics.reduce((sum, e) => sum + e.requestCount, 0);
		const totalErrors = mockEndpointMetrics.reduce((sum, e) => sum + e.errorCount, 0);
		const avgLatency =
			mockEndpointMetrics.reduce((sum, e) => sum + e.avgDuration, 0) / mockEndpointMetrics.length;
		const p95Latency =
			mockEndpointMetrics.reduce((sum, e) => sum + e.p95, 0) / mockEndpointMetrics.length;
		const slowTraces = mockTraces.filter((t) => t.duration > 1000).length;
		const errorTraces = mockTraces.filter((t) => t.statusCode >= 400).length;
		const slowQueries = mockSlowQueries.filter((q) => q.duration > 1000).length;

		return {
			totalRequests,
			totalErrors,
			errorRate: ((totalErrors / totalRequests) * 100).toFixed(2),
			avgLatency: Math.round(avgLatency),
			p95Latency: Math.round(p95Latency),
			slowTraces,
			errorTraces,
			slowQueries
		};
	});

	function toggleTrace(traceId: string) {
		expandedTrace = expandedTrace === traceId ? null : traceId;
	}

	function toggleEndpoint(endpoint: string) {
		expandedEndpoint = expandedEndpoint === endpoint ? null : endpoint;
	}

	function toggleQuery(queryId: string) {
		expandedQuery = expandedQuery === queryId ? null : queryId;
	}

	function formatDuration(ms: number): string {
		if (ms < 1000) return `${ms}ms`;
		return `${(ms / 1000).toFixed(2)}s`;
	}

	function formatTimestamp(date: Date): string {
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);

		if (diffMins < 1) return 'Just now';
		if (diffMins < 60) return `${diffMins}m ago`;
		if (diffMins < 1440) return `${Math.floor(diffMins / 60)}h ago`;
		return `${Math.floor(diffMins / 1440)}d ago`;
	}

	function getDurationColor(ms: number): string {
		if (ms < 200) return 'text-green-600';
		if (ms < 1000) return 'text-yellow-600';
		return 'text-red-600';
	}

	function getStatusColor(code: number): string {
		if (code >= 200 && code < 300) return 'bg-green-100 text-green-800';
		if (code >= 400 && code < 500) return 'bg-yellow-100 text-yellow-800';
		return 'bg-red-100 text-red-800';
	}

	function getSpanTypeIcon(type: string) {
		switch (type) {
			case 'database':
				return Database;
			case 'cache':
				return Zap;
			case 'external':
				return TrendingUp;
			default:
				return Activity;
		}
	}

	function getSpanColor(status: string): string {
		switch (status) {
			case 'success':
				return 'bg-green-500';
			case 'warning':
				return 'bg-yellow-500';
			case 'error':
				return 'bg-red-500';
			default:
				return 'bg-gray-500';
		}
	}

	function calculateSpanPosition(span: TraceSpan, totalDuration: number) {
		const left = (span.startOffset / totalDuration) * 100;
		const width = (span.duration / totalDuration) * 100;
		return { left: `${left}%`, width: `${width}%` };
	}

	async function handleRefresh() {
		isRefreshing = true;
		await new Promise((resolve) => setTimeout(resolve, 1000));
		isRefreshing = false;
	}

	function exportToJSON() {
		const data = {
			traces: mockTraces,
			endpoints: mockEndpointMetrics,
			slowQueries: mockSlowQueries,
			stats: stats(),
			exportedAt: new Date().toISOString()
		};
		const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = `apm-data-${Date.now()}.json`;
		a.click();
		URL.revokeObjectURL(url);
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold">Application Performance Monitoring</h1>
			<p class="text-muted-foreground">
				Monitor request traces, latency metrics, and database performance
			</p>
		</div>
		<div class="flex gap-2">
			<Button variant="outline" onclick={exportToJSON}>
				<Download class="mr-2 h-4 w-4" />
				Export
			</Button>
		<Button variant="outline" onclick={handleRefresh} disabled={isRefreshing}>
			<RefreshCw class={`mr-2 h-4 w-4 ${isRefreshing ? 'animate-spin' : ''}`} />
			Refresh
		</Button>
		</div>
	</div>

	<div class="grid gap-4 md:grid-cols-4 lg:grid-cols-8">
		<div class="rounded-lg border bg-card p-6 md:col-span-2 lg:col-span-2">
			<div class="flex items-center gap-2">
				<Activity class="h-4 w-4 text-muted-foreground" />
				<p class="text-sm font-medium text-muted-foreground">Total Requests</p>
			</div>
			<p class="mt-2 text-2xl font-bold">{stats().totalRequests.toLocaleString()}</p>
		</div>

		<div class="rounded-lg border bg-card p-6 md:col-span-2 lg:col-span-2">
			<div class="flex items-center gap-2">
				<AlertTriangle class="h-4 w-4 text-muted-foreground" />
				<p class="text-sm font-medium text-muted-foreground">Error Rate</p>
			</div>
			<p class="mt-2 text-2xl font-bold">{stats().errorRate}%</p>
			<p class="text-xs text-muted-foreground">{stats().totalErrors} errors</p>
		</div>

		<div class="rounded-lg border bg-card p-6 md:col-span-2 lg:col-span-2">
			<div class="flex items-center gap-2">
				<Clock class="h-4 w-4 text-muted-foreground" />
				<p class="text-sm font-medium text-muted-foreground">Avg Latency</p>
			</div>
			<p class="mt-2 text-2xl font-bold">{stats().avgLatency}ms</p>
			<p class="text-xs text-muted-foreground">p95: {stats().p95Latency}ms</p>
		</div>

		<div class="rounded-lg border bg-card p-6 md:col-span-2 lg:col-span-2">
			<div class="flex items-center gap-2">
				<TrendingUp class="h-4 w-4 text-muted-foreground" />
				<p class="text-sm font-medium text-muted-foreground">Slow Traces</p>
			</div>
			<p class="mt-2 text-2xl font-bold">{stats().slowTraces}</p>
			<p class="text-xs text-muted-foreground">&gt;1s response time</p>
		</div>

		<div class="rounded-lg border bg-card p-6 md:col-span-2 lg:col-span-2">
			<div class="flex items-center gap-2">
				<AlertTriangle class="h-4 w-4 text-muted-foreground" />
				<p class="text-sm font-medium text-muted-foreground">Error Traces</p>
			</div>
			<p class="mt-2 text-2xl font-bold">{stats().errorTraces}</p>
		</div>

		<div class="rounded-lg border bg-card p-6 md:col-span-2 lg:col-span-2">
			<div class="flex items-center gap-2">
				<Database class="h-4 w-4 text-muted-foreground" />
				<p class="text-sm font-medium text-muted-foreground">Slow Queries</p>
			</div>
			<p class="mt-2 text-2xl font-bold">{stats().slowQueries}</p>
			<p class="text-xs text-muted-foreground">&gt;1s query time</p>
		</div>

		<div class="rounded-lg border bg-card p-6 md:col-span-2 lg:col-span-4">
			<div class="flex items-center gap-2">
				<Zap class="h-4 w-4 text-muted-foreground" />
				<p class="text-sm font-medium text-muted-foreground">Performance Score</p>
			</div>
			<p class="mt-2 text-2xl font-bold text-green-600">87.4</p>
			<p class="text-xs text-muted-foreground">Good performance</p>
		</div>
	</div>

	<div class="grid gap-6 md:grid-cols-2">
		<div class="rounded-lg border bg-card p-6">
			<h3 class="mb-4 text-sm font-medium">Latency Over Time</h3>
			<LineChart data={latencyChartData} unit="ms" height={250} color="#3b82f6" />
		</div>

		<div class="rounded-lg border bg-card p-6">
			<h3 class="mb-4 text-sm font-medium">Request Rate</h3>
			<AreaChart data={requestRateChartData} unit="req/h" height={250} color="#10b981" />
		</div>

		<div class="rounded-lg border bg-card p-6">
			<h3 class="mb-4 text-sm font-medium">Error Rate Over Time</h3>
			<LineChart data={errorRateChartData} unit="%" height={250} color="#ef4444" />
		</div>

		<div class="rounded-lg border bg-card p-6">
			<h3 class="mb-4 text-sm font-medium">Throughput by Method</h3>
			<BarChart data={throughputChartData} unit="requests" height={250} color="#8b5cf6" />
		</div>
	</div>

	<div class="rounded-lg border bg-card">
		<div class="border-b p-4">
			<div class="flex gap-4">
				<button
					class="px-4 py-2 text-sm font-medium transition-colors {activeTab === 'traces'
						? 'border-b-2 border-primary text-primary'
						: 'text-muted-foreground hover:text-foreground'}"
					onclick={() => (activeTab = 'traces')}
				>
					Request Traces
				</button>
				<button
					class="px-4 py-2 text-sm font-medium transition-colors {activeTab === 'endpoints'
						? 'border-b-2 border-primary text-primary'
						: 'text-muted-foreground hover:text-foreground'}"
					onclick={() => (activeTab = 'endpoints')}
				>
					Endpoint Metrics
				</button>
				<button
					class="px-4 py-2 text-sm font-medium transition-colors {activeTab === 'queries'
						? 'border-b-2 border-primary text-primary'
						: 'text-muted-foreground hover:text-foreground'}"
					onclick={() => (activeTab = 'queries')}
				>
					Slow Queries
				</button>
			</div>
		</div>

		<div class="p-6">
			<div class="mb-6 flex flex-col gap-4 md:flex-row">
				<div class="relative flex-1">
					<Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
					<Input
						type="text"
						placeholder="Search {activeTab}..."
						class="pl-10"
						bind:value={searchQuery}
					/>
				</div>
				{#if activeTab === 'traces' || activeTab === 'endpoints'}
					<Select bind:selected={statusFilter}>
						<SelectTrigger class="w-full md:w-48">
							<SelectValue placeholder="Filter by status" />
						</SelectTrigger>
						<SelectContent>
							<SelectItem value="all">All Status</SelectItem>
							<SelectItem value="success">Success</SelectItem>
							<SelectItem value="error">Errors</SelectItem>
							<SelectItem value="slow">Slow (&gt;1s)</SelectItem>
						</SelectContent>
					</Select>
				{/if}
				{#if activeTab === 'traces' || activeTab === 'endpoints'}
					<Select bind:selected={methodFilter}>
						<SelectTrigger class="w-full md:w-48">
							<SelectValue placeholder="Filter by method" />
						</SelectTrigger>
						<SelectContent>
							<SelectItem value="all">All Methods</SelectItem>
							<SelectItem value="GET">GET</SelectItem>
							<SelectItem value="POST">POST</SelectItem>
							<SelectItem value="PUT">PUT</SelectItem>
							<SelectItem value="DELETE">DELETE</SelectItem>
						</SelectContent>
					</Select>
				{/if}
			</div>

			{#if activeTab === 'traces'}
				<div class="space-y-4">
					{#each filteredTraces() as trace (trace.id)}
						<div class="rounded-lg border bg-card">
							<div class="p-4">
								<div class="flex items-start justify-between">
									<div class="flex-1">
										<div class="flex items-center gap-3">
											<span class="rounded px-2 py-1 text-xs font-medium {getStatusColor(trace.statusCode)}">
												{trace.statusCode}
											</span>
											<span class="rounded bg-secondary px-2 py-1 text-xs font-medium">
												{trace.method}
											</span>
											<span class="font-mono text-sm font-medium">{trace.endpoint}</span>
											<span class="text-sm font-medium {getDurationColor(trace.duration)}">
												{formatDuration(trace.duration)}
											</span>
										</div>
										<div class="mt-2 flex items-center gap-4 text-sm text-muted-foreground">
											<span>Trace ID: {trace.id}</span>
											<span>{formatTimestamp(trace.timestamp)}</span>
											<span>{trace.spans.length} spans</span>
										</div>
										{#if trace.errorMessage}
											<div
												class="mt-2 flex items-center gap-2 rounded bg-red-50 p-2 text-sm text-red-800"
											>
												<AlertTriangle class="h-4 w-4" />
												<span>{trace.errorMessage}</span>
											</div>
										{/if}
									</div>
									<Button
										variant="ghost"
										size="sm"
										onclick={() => toggleTrace(trace.id)}
									>
										{#if expandedTrace === trace.id}
											<ChevronUp class="h-4 w-4" />
										{:else}
											<ChevronDown class="h-4 w-4" />
										{/if}
									</Button>
								</div>

								{#if expandedTrace === trace.id}
									<div class="mt-4 space-y-4 border-t pt-4">
										<div>
											<h4 class="mb-3 text-sm font-medium">Trace Timeline</h4>
											<div class="space-y-2">
												{#each trace.spans as span (span.id)}
													<div class="flex items-center gap-3">
														<div class="flex w-32 items-center gap-2">
															<svelte:component
																this={getSpanTypeIcon(span.type)}
																class="h-4 w-4 text-muted-foreground"
															/>
															<span class="text-xs text-muted-foreground">{span.type}</span>
														</div>
														<div class="flex-1">
															<div class="relative h-8 rounded-md bg-secondary">
																<div
																	class="absolute top-0 h-full rounded-md {getSpanColor(span.status)}"
																	style="{calculateSpanPosition(span, trace.duration).left}; width: {calculateSpanPosition(span, trace.duration).width}"
																></div>
																<div class="relative flex h-full items-center px-2">
																	<span class="text-xs font-medium text-white mix-blend-difference">
																		{span.name} - {formatDuration(span.duration)}
																	</span>
																</div>
															</div>
															{#if span.details}
																<p class="mt-1 truncate text-xs text-muted-foreground">
																	{span.details}
																</p>
															{/if}
														</div>
													</div>
												{/each}
											</div>
										</div>
									</div>
								{/if}
							</div>
						</div>
					{/each}

					{#if filteredTraces().length === 0}
						<div class="rounded-lg border border-dashed p-8 text-center">
							<Activity class="mx-auto h-12 w-12 text-muted-foreground opacity-50" />
							<p class="mt-4 text-sm font-medium">No traces found</p>
							<p class="text-sm text-muted-foreground">
								Try adjusting your search or filter criteria
							</p>
						</div>
					{/if}
				</div>
			{:else if activeTab === 'endpoints'}
				<div class="space-y-4">
					{#each filteredEndpoints() as endpoint (endpoint.endpoint + endpoint.method)}
						<div class="rounded-lg border bg-card">
							<div class="p-4">
								<div class="flex items-start justify-between">
									<div class="flex-1">
										<div class="flex items-center gap-3">
											<span class="rounded bg-secondary px-2 py-1 text-xs font-medium">
												{endpoint.method}
											</span>
											<span class="font-mono text-sm font-medium">{endpoint.endpoint}</span>
											<span class="text-sm font-medium {getDurationColor(endpoint.avgDuration)}">
												avg: {formatDuration(endpoint.avgDuration)}
											</span>
										</div>
										<div class="mt-2 grid grid-cols-4 gap-4 text-sm">
											<div>
												<p class="text-muted-foreground">Requests</p>
												<p class="font-medium">{endpoint.requestCount.toLocaleString()}</p>
											</div>
											<div>
												<p class="text-muted-foreground">Error Rate</p>
												<p class="font-medium {endpoint.errorRate > 5 ? 'text-red-600' : endpoint.errorRate > 1 ? 'text-yellow-600' : 'text-green-600'}">
													{endpoint.errorRate.toFixed(2)}%
												</p>
											</div>
											<div>
												<p class="text-muted-foreground">p95</p>
												<p class="font-medium {getDurationColor(endpoint.p95)}">
													{formatDuration(endpoint.p95)}
												</p>
											</div>
											<div>
												<p class="text-muted-foreground">p99</p>
												<p class="font-medium {getDurationColor(endpoint.p99)}">
													{formatDuration(endpoint.p99)}
												</p>
											</div>
										</div>
									</div>
									<Button
										variant="ghost"
										size="sm"
										onclick={() => toggleEndpoint(endpoint.endpoint)}
									>
										{#if expandedEndpoint === endpoint.endpoint}
											<ChevronUp class="h-4 w-4" />
										{:else}
											<ChevronDown class="h-4 w-4" />
										{/if}
									</Button>
								</div>

								{#if expandedEndpoint === endpoint.endpoint}
									<div class="mt-4 space-y-3 border-t pt-4">
										<div class="grid grid-cols-3 gap-4 text-sm">
											<div>
												<p class="text-muted-foreground">p50 Latency</p>
												<p class="font-medium">{formatDuration(endpoint.p50)}</p>
											</div>
											<div>
												<p class="text-muted-foreground">Total Errors</p>
												<p class="font-medium">{endpoint.errorCount}</p>
											</div>
											<div>
												<p class="text-muted-foreground">Slowest Trace</p>
												<p class="font-mono text-xs">{endpoint.slowestTrace || 'N/A'}</p>
											</div>
										</div>
									</div>
								{/if}
							</div>
						</div>
					{/each}

					{#if filteredEndpoints().length === 0}
						<div class="rounded-lg border border-dashed p-8 text-center">
							<TrendingUp class="mx-auto h-12 w-12 text-muted-foreground opacity-50" />
							<p class="mt-4 text-sm font-medium">No endpoints found</p>
							<p class="text-sm text-muted-foreground">
								Try adjusting your search or filter criteria
							</p>
						</div>
					{/if}
				</div>
			{:else if activeTab === 'queries'}
				<div class="space-y-4">
					{#each filteredQueries() as query (query.id)}
						<div class="rounded-lg border bg-card">
							<div class="p-4">
								<div class="flex items-start justify-between">
									<div class="flex-1">
										<div class="flex items-center gap-3">
											<Database class="h-4 w-4 text-muted-foreground" />
											<span class="text-sm font-medium">{query.database}</span>
											<span class="text-sm font-medium {getDurationColor(query.duration)}">
												{formatDuration(query.duration)}
											</span>
										</div>
										<div class="mt-2 rounded bg-secondary p-3">
											<code class="text-xs">{query.query}</code>
										</div>
										<div class="mt-2 flex items-center gap-4 text-sm text-muted-foreground">
											<span>{formatTimestamp(query.timestamp)}</span>
											<span>Executed {query.executionCount}x</span>
											<span>{query.affectedRows} rows</span>
											<span>Trace: {query.traceId}</span>
										</div>
									</div>
									<Button
										variant="ghost"
										size="sm"
										onclick={() => toggleQuery(query.id)}
									>
										{#if expandedQuery === query.id}
											<ChevronUp class="h-4 w-4" />
										{:else}
											<ChevronDown class="h-4 w-4" />
										{/if}
									</Button>
								</div>

								{#if expandedQuery === query.id}
									<div class="mt-4 space-y-3 border-t pt-4">
										<div>
											<h4 class="text-sm font-medium">Query Details</h4>
											<div class="mt-2 grid grid-cols-3 gap-4 text-sm">
												<div>
													<p class="text-muted-foreground">Query ID</p>
													<p class="font-mono text-xs">{query.id}</p>
												</div>
												<div>
													<p class="text-muted-foreground">Database</p>
													<p class="font-medium">{query.database}</p>
												</div>
												<div>
													<p class="text-muted-foreground">Duration</p>
													<p class="font-medium">{formatDuration(query.duration)}</p>
												</div>
												<div>
													<p class="text-muted-foreground">Execution Count</p>
													<p class="font-medium">{query.executionCount}</p>
												</div>
												<div>
													<p class="text-muted-foreground">Affected Rows</p>
													<p class="font-medium">{query.affectedRows}</p>
												</div>
												<div>
													<p class="text-muted-foreground">Related Trace</p>
													<p class="font-mono text-xs">{query.traceId}</p>
												</div>
											</div>
										</div>
										<div class="rounded bg-yellow-50 p-3">
											<p class="text-sm text-yellow-800">
												<strong>Optimization Suggestion:</strong>
												{#if query.executionCount > 100}
													Consider caching this frequently executed query.
												{:else if query.affectedRows > 1000}
													Large result set - consider pagination or indexing.
												{:else}
													Review query execution plan and add appropriate indexes.
												{/if}
											</p>
										</div>
									</div>
								{/if}
							</div>
						</div>
					{/each}

					{#if filteredQueries().length === 0}
						<div class="rounded-lg border border-dashed p-8 text-center">
							<Database class="mx-auto h-12 w-12 text-muted-foreground opacity-50" />
							<p class="mt-4 text-sm font-medium">No slow queries found</p>
							<p class="text-sm text-muted-foreground">
								Try adjusting your search criteria
							</p>
						</div>
					{/if}
				</div>
			{/if}
		</div>
	</div>
</div>
