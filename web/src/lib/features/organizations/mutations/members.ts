import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import * as membersApi from '$lib/api/organizations/members';
import { membersKeys } from '../keys/members';
import type { OrganizationMember, InviteMemberRequest, UpdateMemberRoleRequest } from '$lib/api/organizations/members';

type InviteMemberMutationOptions = {
	onSuccess?: () => void;
	onError?: (error: Error) => void;
};

export const inviteMemberMutation = (organizationId: string, options: InviteMemberMutationOptions = {}) => {
	const queryClient = useQueryClient();

	return createMutation<{ message: string }, Error, InviteMemberRequest>(() => ({
		mutationFn: (data) => membersApi.inviteMember(organizationId, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: membersKeys.list(organizationId) });
			options.onSuccess?.();
		},
		onError: (error: Error) => {
			console.error(error);
			options.onError?.(error);
		}
	}));
};

type UpdateMemberRoleMutationOptions = {
	onSuccess?: () => void;
	onError?: (error: Error) => void;
};

export const updateMemberRoleMutation = (organizationId: string, options: UpdateMemberRoleMutationOptions = {}) => {
	const queryClient = useQueryClient();

	return createMutation<
		{ message: string },
		Error,
		{ memberId: string; data: UpdateMemberRoleRequest }
	>(() => ({
		mutationFn: ({ memberId, data }) => membersApi.updateMemberRole(organizationId, memberId, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: membersKeys.list(organizationId) });
			options.onSuccess?.();
		},
		onError: (error: Error) => {
			console.error(error);
			options.onError?.(error);
		}
	}));
};

type RemoveMemberMutationOptions = {
	onSuccess?: () => void;
	onError?: (error: Error) => void;
};

export const removeMemberMutation = (organizationId: string, options: RemoveMemberMutationOptions = {}) => {
	const queryClient = useQueryClient();

	return createMutation<void, Error, string>(() => ({
		mutationFn: (memberId) => membersApi.removeMember(organizationId, memberId),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: membersKeys.list(organizationId) });
			options.onSuccess?.();
		},
		onError: (error: Error) => {
			console.error(error);
			options.onError?.(error);
		}
	}));
};
