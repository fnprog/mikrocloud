export interface Organization {
	id: string;
	name: string;
	slug: string;
	description: string;
	owner_id: string;
	billing_email: string;
	plan: 'free' | 'pro' | 'enterprise';
	status: 'active' | 'suspended' | 'deleted';
	created_at: string;
	updated_at: string;
}
