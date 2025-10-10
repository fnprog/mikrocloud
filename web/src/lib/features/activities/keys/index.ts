export const activitiesKeys = {
  all: ['activities'] as const,
  recent: (limit?: number, offset?: number) =>
    [...activitiesKeys.all, 'recent', limit, offset] as const,
  forResource: (resourceType: string, resourceId: string, limit?: number, offset?: number) =>
    [...activitiesKeys.all, 'resource', resourceType, resourceId, limit, offset] as const
};
