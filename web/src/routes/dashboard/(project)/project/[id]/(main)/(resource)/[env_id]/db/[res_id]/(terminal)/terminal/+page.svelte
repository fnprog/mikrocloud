<script>
	import { page } from '$app/stores';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import {
		Terminal,
		Play,
		Square,
		RotateCcw,
		Download,
		Upload,
		Settings,
		Maximize2,
		Minimize2
	} from 'lucide-svelte';

	let projectId = $derived($page.params.id);
	let envId = $derived($page.params.env_id);
	let resId = $derived($page.params.res_id);

	let isConnected = $state(false);
	let isFullscreen = $state(false);
	let currentCommand = $state('');
	let terminalHistory = $state([
		{ type: 'system', content: 'Welcome to Database Terminal' },
		{ type: 'system', content: 'Type "help" for available commands' },
		{ type: 'prompt', content: 'user@database:~$ ' }
	]);

	let sessions = $state([
		{
			id: 1,
			name: 'Main Session',
			status: 'active',
			startTime: '14:30:25',
			duration: '00:45:12'
		},
		{
			id: 2,
			name: 'Debug Session',
			status: 'idle',
			startTime: '13:15:10',
			duration: '01:20:33'
		}
	]);

	let activeSession = $state(1);

	function connectTerminal() {
		isConnected = true;
		terminalHistory = [
			...terminalHistory,
			{ type: 'system', content: 'Connected to database container' },
			{ type: 'prompt', content: 'root@database:~$ ' }
		];
	}

	function disconnectTerminal() {
		isConnected = false;
		terminalHistory = [
			...terminalHistory,
			{ type: 'system', content: 'Connection closed' },
			{ type: 'prompt', content: 'user@database:~$ ' }
		];
	}

	function executeCommand() {
		if (!currentCommand.trim()) return;

		terminalHistory = [
			...terminalHistory,
			{
				type: 'command',
				content: `${isConnected ? 'root' : 'user'}@database:~$ ${currentCommand}`
			}
		];

		let response = '';
		switch (currentCommand.toLowerCase().trim()) {
			case 'help':
				response = `Available commands:
  ls          - List directory contents
  pwd         - Print working directory
  ps aux      - List running processes
  df -h       - Show disk usage
  top         - Show system processes
  tail -f     - Follow log files
  psql        - PostgreSQL client
  mysql       - MySQL client
  mongosh     - MongoDB shell
  redis-cli   - Redis client
  clear       - Clear terminal`;
				break;
			case 'ls':
				response = 'data  config  logs  backups  scripts';
				break;
			case 'pwd':
				response = '/var/lib/database';
				break;
			case 'ps aux':
				response = `USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root         1  0.0  0.1  18508  3396 ?        Ss   14:30   0:00 /sbin/init
postgres   123  2.1  5.4 932844 55432 ?        Sl   14:30   0:45 postgres: main
postgres   456  0.0  0.2  12345  2048 ?        S    14:35   0:00 postgres: checkpointer`;
				break;
			case 'df -h':
				response = `Filesystem      Size  Used Avail Use% Mounted on
/dev/sda1        50G   32G   16G  67% /
/dev/sda2       100G   67G   28G  71% /data
tmpfs           2.0G     0  2.0G   0% /dev/shm`;
				break;
			case 'clear':
				terminalHistory = [
					{ type: 'system', content: 'Terminal cleared' },
					{ type: 'prompt', content: `${isConnected ? 'root' : 'user'}@database:~$ ` }
				];
				currentCommand = '';
				return;
			default:
				if (currentCommand.startsWith('tail -f')) {
					response = `Following ${currentCommand.split(' ')[2] || '/var/log/database.log'}...
2024-01-31 14:30:15 [INFO] Database started successfully
2024-01-31 14:30:16 [INFO] Accepting connections on port 5432
2024-01-31 14:30:17 [INFO] Checkpoint completed`;
				} else {
					response = `bash: ${currentCommand}: command not found`;
				}
		}

		terminalHistory = [
			...terminalHistory,
			{ type: 'output', content: response },
			{ type: 'prompt', content: `${isConnected ? 'root' : 'user'}@database:~$ ` }
		];

		currentCommand = '';
	}

	function handleKeydown(event) {
		if (event.key === 'Enter') {
			executeCommand();
		}
	}

	function clearTerminal() {
		terminalHistory = [
			{ type: 'system', content: 'Terminal cleared' },
			{ type: 'prompt', content: `${isConnected ? 'root' : 'user'}@database:~$ ` }
		];
	}

	function toggleFullscreen() {
		isFullscreen = !isFullscreen;
	}

	function switchSession(sessionId) {
		activeSession = sessionId;
		terminalHistory = [
			{ type: 'system', content: `Switched to session ${sessionId}` },
			{ type: 'prompt', content: `${isConnected ? 'root' : 'user'}@database:~$ ` }
		];
	}
</script>

<svelte:head>
	<title>Database Terminal</title>
</svelte:head>

<div class="flex-1 p-6">
	<div class="flex items-center justify-between mb-6">
		<div class="flex items-center space-x-4">
			<div>
				<h1 class="text-2xl font-semibold text-gray-900">Database Terminal</h1>
				<p class="text-sm text-gray-500 mt-1">
					Access your database container directly through the web terminal.
				</p>
			</div>
		</div>
		<div class="flex items-center space-x-3">
			{#if !isConnected}
				<Button onclick={connectTerminal}>
					<Play class="w-4 h-4 mr-2" />
					Connect
				</Button>
			{:else}
				<Button variant="outline" onclick={disconnectTerminal}>
					<Square class="w-4 h-4 mr-2" />
					Disconnect
				</Button>
			{/if}
			<Button variant="outline" onclick={clearTerminal}>
				<RotateCcw class="w-4 h-4 mr-2" />
				Clear
			</Button>
			<Button variant="outline" onclick={toggleFullscreen}>
				{#if isFullscreen}
					<Minimize2 class="w-4 h-4" />
				{:else}
					<Maximize2 class="w-4 h-4" />
				{/if}
			</Button>
		</div>
	</div>

	<div class="mb-4">
		<div class="flex items-center space-x-2">
			{#each sessions as session (session.id)}
				<button
					onclick={() => switchSession(session.id)}
					class="px-3 py-1 text-sm rounded-md border {activeSession === session.id
						? 'bg-blue-50 border-blue-200 text-blue-700'
						: 'bg-white border-gray-200 text-gray-600 hover:bg-gray-50'}"
				>
					<div class="flex items-center space-x-2">
						<div
							class="w-2 h-2 rounded-full {session.status === 'active'
								? 'bg-green-500'
								: 'bg-gray-400'}"
						></div>
						<span>{session.name}</span>
					</div>
				</button>
			{/each}
		</div>
	</div>

	<Card class={isFullscreen ? 'fixed inset-4 z-50' : 'h-[calc(100vh-300px)]'}>
		<CardHeader class="pb-2">
			<div class="flex items-center justify-between">
				<div class="flex items-center space-x-2">
					<Terminal class="w-5 h-5 text-gray-600" />
					<CardTitle class="text-lg">Terminal Session {activeSession}</CardTitle>
					<div class="flex items-center space-x-1">
						<div class="w-2 h-2 rounded-full {isConnected ? 'bg-green-500' : 'bg-red-500'}"></div>
						<span class="text-sm text-gray-600">{isConnected ? 'Connected' : 'Disconnected'}</span>
					</div>
				</div>
				<div class="flex items-center space-x-2">
					<Button size="sm" variant="outline">
						<Download class="w-4 h-4" />
					</Button>
					<Button size="sm" variant="outline">
						<Upload class="w-4 h-4" />
					</Button>
					<Button size="sm" variant="outline">
						<Settings class="w-4 h-4" />
					</Button>
				</div>
			</div>
		</CardHeader>
		<CardContent class="p-0 h-full">
			<div class="h-full bg-black text-green-400 font-mono text-sm overflow-auto">
				<div class="p-4 space-y-1">
					{#each terminalHistory as line}
						<div class="flex">
							{#if line.type === 'system'}
								<span class="text-blue-400"># {line.content}</span>
							{:else if line.type === 'command'}
								<span class="text-white">{line.content}</span>
							{:else if line.type === 'output'}
								<pre class="text-gray-300 whitespace-pre-wrap">{line.content}</pre>
							{:else if line.type === 'prompt'}
								<span class="text-green-400">{line.content}</span>
							{/if}
						</div>
					{/each}

					<div class="flex items-center">
						<span class="text-green-400"
							>{isConnected ? 'root' : 'user'}@database:~$
						</span>
						<input
							bind:value={currentCommand}
							onkeydown={handleKeydown}
							class="flex-1 bg-transparent text-white outline-none ml-1"
							placeholder="Type a command..."
							autocomplete="off"
							spellcheck="false"
						/>
						<span class="animate-pulse text-white">|</span>
					</div>
				</div>
			</div>
		</CardContent>
	</Card>
</div>
