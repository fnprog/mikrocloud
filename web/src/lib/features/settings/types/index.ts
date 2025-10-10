export interface GeneralSettings {
	domain: string;
	timezone: string;
	ipv4: string;
	ipv6: string;
	allow_registrations: boolean;
	do_not_track: boolean;
}

export interface UpdateGeneralSettingsRequest {
	domain: string;
	timezone: string;
	ipv4: string;
	ipv6: string;
	allow_registrations: boolean;
	do_not_track: boolean;
}

export interface DetectedIPs {
	ipv4: string;
	ipv6: string;
}
