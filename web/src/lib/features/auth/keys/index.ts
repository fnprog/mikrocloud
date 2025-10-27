export const authKeys = {
	all: ['auth'] as const,
	profile: () => [...authKeys.all, 'profile'] as const,
	setupStatus: () => [...authKeys.all, 'setup-status'] as const
};
