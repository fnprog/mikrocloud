import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { settingsApi } from '../api';
import { settingsKeys } from '../keys';
import type { UpdateGeneralSettingsRequest } from "../types"

export const createUpdateGeneralSettingsMutation = () => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: (data: UpdateGeneralSettingsRequest) => settingsApi.updateGeneral(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: settingsKeys.general() });
    }
  }));
};
