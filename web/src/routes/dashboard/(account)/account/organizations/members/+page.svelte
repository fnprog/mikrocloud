<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Select, SelectContent, SelectItem, SelectTrigger } from '$lib/components/ui/select';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { UserPlus, Mail, MoreVertical, Trash2, Crown, Shield } from 'lucide-svelte';

	type MemberRole = 'owner' | 'admin' | 'developer' | 'member' | 'viewer';
	type MemberStatus = 'active' | 'pending' | 'inactive';

	type Member = {
		id: string;
		name: string;
		email: string;
		role: MemberRole;
		status: MemberStatus;
		joinedAt: string;
	};

	let members = $state<Member[]>([
		{
			id: '1',
			name: 'Admin User',
			email: 'admin@example.com',
			role: 'owner',
			status: 'active',
			joinedAt: '2024-01-15'
		},
		{
			id: '2',
			name: 'John Doe',
			email: 'john@example.com',
			role: 'admin',
			status: 'active',
			joinedAt: '2024-02-20'
		},
		{
			id: '3',
			name: 'Jane Smith',
			email: 'jane@example.com',
			role: 'developer',
			status: 'active',
			joinedAt: '2024-03-10'
		},
		{
			id: '4',
			name: 'Bob Wilson',
			email: 'bob@example.com',
			role: 'member',
			status: 'pending',
			joinedAt: '2024-04-01'
		}
	]);

	let isInviteDialogOpen = $state(false);
	let isRemoveDialogOpen = $state(false);
	let selectedMember = $state<Member | null>(null);

	let inviteEmail = $state('');
	let inviteRole = $state<MemberRole>('member');

	function getRoleBadgeVariant(role: MemberRole) {
		switch (role) {
			case 'owner':
				return 'default';
			case 'admin':
				return 'secondary';
			case 'developer':
				return 'outline';
			case 'member':
				return 'outline';
			case 'viewer':
				return 'outline';
		}
	}

	function getStatusBadgeVariant(status: MemberStatus) {
		switch (status) {
			case 'active':
				return 'default';
			case 'pending':
				return 'secondary';
			case 'inactive':
				return 'outline';
		}
	}

	function handleInvite() {
		console.log('Inviting:', inviteEmail, 'as', inviteRole);
		inviteEmail = '';
		inviteRole = 'member';
		isInviteDialogOpen = false;
	}

	function handleRemoveMember() {
		if (selectedMember) {
			console.log('Removing member:', selectedMember.id);
			members = members.filter((m) => m.id !== selectedMember.id);
		}
		isRemoveDialogOpen = false;
		selectedMember = null;
	}

	function openRemoveDialog(member: Member) {
		selectedMember = member;
		isRemoveDialogOpen = true;
	}

	function handleRoleChange(memberId: string, newRole: MemberRole) {
		const memberIndex = members.findIndex((m) => m.id === memberId);
		if (memberIndex !== -1) {
			members[memberIndex].role = newRole;
		}
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold mb-2">Members</h1>
			<p class="text-muted-foreground">
				Manage your organization's members and their permissions.
			</p>
		</div>
		<Button onclick={() => (isInviteDialogOpen = true)}>
			<UserPlus class="w-4 h-4 mr-2" />
			Invite Member
		</Button>
	</div>

	<Card>
		<CardHeader>
			<CardTitle>Organization Members ({members.length})</CardTitle>
			<CardDescription>View and manage all members in your organization</CardDescription>
		</CardHeader>
		<CardContent>
			<div class="space-y-4">
				{#each members as member (member.id)}
					<div class="flex items-center justify-between p-4 border rounded-lg">
						<div class="flex items-center gap-4 flex-1">
							<div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
								<span class="text-sm font-medium">{member.name.charAt(0)}</span>
							</div>
							<div class="flex-1">
								<div class="flex items-center gap-2">
									<p class="font-medium">{member.name}</p>
									{#if member.role === 'owner'}
										<Crown class="w-4 h-4 text-yellow-500" />
									{:else if member.role === 'admin'}
										<Shield class="w-4 h-4 text-blue-500" />
									{/if}
								</div>
								<div class="flex items-center gap-2 text-sm text-muted-foreground">
									<Mail class="w-3 h-3" />
									<span>{member.email}</span>
								</div>
							</div>
						</div>
						<div class="flex items-center gap-3">
							<Badge variant={getStatusBadgeVariant(member.status)}>
								{member.status}
							</Badge>
							<div class="w-32">
								<Select
									type="single"
									value={member.role}
									disabled={member.role === 'owner'}
									onValueChange={(value) => value && handleRoleChange(member.id, value as MemberRole)}
								>
									<SelectTrigger>
										{member.role}
									</SelectTrigger>
									<SelectContent>
										<SelectItem value="admin">Admin</SelectItem>
										<SelectItem value="developer">Developer</SelectItem>
										<SelectItem value="member">Member</SelectItem>
										<SelectItem value="viewer">Viewer</SelectItem>
									</SelectContent>
								</Select>
							</div>
							{#if member.role !== 'owner'}
								<Button
									variant="ghost"
									size="icon"
									onclick={() => openRemoveDialog(member)}
								>
									<Trash2 class="w-4 h-4 text-destructive" />
								</Button>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		</CardContent>
	</Card>

	<Card>
		<CardHeader>
			<CardTitle>Roles & Permissions</CardTitle>
			<CardDescription>Understanding member roles in your organization</CardDescription>
		</CardHeader>
		<CardContent>
			<div class="space-y-3">
				<div class="flex items-start gap-3 text-sm">
					<Badge variant="default" class="mt-0.5">Owner</Badge>
					<p class="text-muted-foreground">Full access to all organization settings, billing, and member management. Can delete the organization.</p>
				</div>
				<div class="flex items-start gap-3 text-sm">
					<Badge variant="secondary" class="mt-0.5">Admin</Badge>
					<p class="text-muted-foreground">Can manage projects, applications, and invite members. Cannot manage billing or delete the organization.</p>
				</div>
				<div class="flex items-start gap-3 text-sm">
					<Badge variant="outline" class="mt-0.5">Developer</Badge>
					<p class="text-muted-foreground">Can create and manage projects, applications, and deployments. Cannot manage members.</p>
				</div>
				<div class="flex items-start gap-3 text-sm">
					<Badge variant="outline" class="mt-0.5">Member</Badge>
					<p class="text-muted-foreground">Can view projects and create deployments. Limited modification access.</p>
				</div>
				<div class="flex items-start gap-3 text-sm">
					<Badge variant="outline" class="mt-0.5">Viewer</Badge>
					<p class="text-muted-foreground">Read-only access to projects and applications. Cannot make any changes.</p>
				</div>
			</div>
		</CardContent>
	</Card>
</div>

<Dialog.Root bind:open={isInviteDialogOpen}>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Invite Member</Dialog.Title>
			<Dialog.Description>
				Send an invitation to join your organization. They will receive an email with instructions.
			</Dialog.Description>
		</Dialog.Header>
		<div class="space-y-4 py-4">
			<div class="space-y-2">
				<Label for="invite-email">Email Address</Label>
				<Input
					id="invite-email"
					type="email"
					bind:value={inviteEmail}
					placeholder="user@example.com"
				/>
			</div>
			<div class="space-y-2">
				<Label for="invite-role">Role</Label>
				<Select type="single" bind:value={inviteRole}>
					<SelectTrigger id="invite-role">
						{inviteRole}
					</SelectTrigger>
					<SelectContent>
						<SelectItem value="admin">Admin</SelectItem>
						<SelectItem value="developer">Developer</SelectItem>
						<SelectItem value="member">Member</SelectItem>
						<SelectItem value="viewer">Viewer</SelectItem>
					</SelectContent>
				</Select>
			</div>
		</div>
		<Dialog.Footer>
			<Button variant="outline" onclick={() => (isInviteDialogOpen = false)}>Cancel</Button>
			<Button onclick={handleInvite}>Send Invitation</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<AlertDialog.Root bind:open={isRemoveDialogOpen}>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Remove Member</AlertDialog.Title>
			<AlertDialog.Description>
				Are you sure you want to remove {selectedMember?.name} from the organization?
				This action cannot be undone.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
			<AlertDialog.Action onclick={handleRemoveMember}>Remove</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
