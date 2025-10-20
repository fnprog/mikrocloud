import { apiClient } from '../client';

export type MemberRole = 'owner' | 'admin' | 'developer' | 'member' | 'viewer';
export type MemberStatus = 'active' | 'pending' | 'inactive';

export interface OrganizationMember {
	id: string;
	organization_id: string;
	user_id: string;
	user_name: string;
	user_email: string;
	role: MemberRole;
	status: MemberStatus;
	invited_by?: string;
	invited_at?: string;
	joined_at?: string;
}

export interface InviteMemberRequest {
	email: string;
	role: MemberRole;
}

export interface UpdateMemberRoleRequest {
	role: MemberRole;
}

export interface TransferOwnershipRequest {
	new_owner_id: string;
}

export async function listMembers(organizationId: string): Promise<OrganizationMember[]> {
	return apiClient.get<OrganizationMember[]>(`/organizations/${organizationId}/members`);
}

export async function inviteMember(organizationId: string, data: InviteMemberRequest): Promise<{ message: string }> {
	return apiClient.post<{ message: string }>(`/organizations/${organizationId}/members`, data);
}

export async function updateMemberRole(organizationId: string, memberId: string, data: UpdateMemberRoleRequest): Promise<{ message: string }> {
	return apiClient.put<{ message: string }>(`/organizations/${organizationId}/members/${memberId}`, data);
}

export async function removeMember(organizationId: string, memberId: string): Promise<void> {
	return apiClient.delete<void>(`/organizations/${organizationId}/members/${memberId}`);
}

export async function transferOwnership(organizationId: string, data: TransferOwnershipRequest): Promise<void> {
	return apiClient.post<void>(`/organizations/${organizationId}/transfer-ownership`, data);
}
