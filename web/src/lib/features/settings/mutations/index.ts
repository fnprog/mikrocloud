import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { settingsApi } from '../api';
import { settingsKeys } from '../keys';
import type { UpdateGeneralSettingsRequest, UpdateSMTPSettingsRequest } from "../types"

export const createUpdateGeneralSettingsMutation = () => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: (data: UpdateGeneralSettingsRequest) => settingsApi.updateGeneral(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: settingsKeys.general() });
    }
  }));
};

export const createUpdateSMTPSettingsMutation = () => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: (data: UpdateSMTPSettingsRequest) => settingsApi.updateSMTP(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: settingsKeys.smtp() });
    }
  }));
};
