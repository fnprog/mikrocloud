export const templatesKeys = {
	all: ['templates'] as const,
	lists: () => [...templatesKeys.all, 'list'] as const,
	list: (category?: string) => [...templatesKeys.lists(), { category }] as const,
	details: () => [...templatesKeys.all, 'detail'] as const,
	detail: (templateId: string) => [...templatesKeys.details(), { templateId }] as const
};
