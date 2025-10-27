import { createQuery } from '@tanstack/svelte-query';
import { settingsApi } from '../api';
import { settingsKeys } from '../keys';

export const createGeneralSettingsQuery = () =>
	createQuery(() => ({
		queryKey: settingsKeys.general(),
		queryFn: () => settingsApi.getGeneral()
	}));

export const createSMTPSettingsQuery = () =>
	createQuery(() => ({
		queryKey: settingsKeys.smtp(),
		queryFn: () => settingsApi.getSMTP()
	}));
