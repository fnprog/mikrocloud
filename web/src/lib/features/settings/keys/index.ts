export const settingsKeys = {
	all: ['settings'] as const,
	general: () => [...settingsKeys.all, 'general'] as const,
	smtp: () => [...settingsKeys.all, 'smtp'] as const
};
