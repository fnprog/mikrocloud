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

export interface SMTPSettings {
	enabled: boolean;
	host: string;
	port: number;
	username: string;
	password: string;
	from_email: string;
	from_name: string;
}

export interface UpdateSMTPSettingsRequest {
	enabled: boolean;
	host: string;
	port: number;
	username: string;
	password: string;
	from_email: string;
	from_name: string;
}

export interface DetectedIPs {
	ipv4: string;
	ipv6: string;
}
