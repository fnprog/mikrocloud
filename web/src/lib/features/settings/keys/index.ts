export const settingsKeys = {
	all: ['settings'] as const,
	general: () => [...settingsKeys.all, 'general'] as const
};
