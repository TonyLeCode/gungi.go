// import { AuthApiError } from '@supabase/supabase-js';
import type { PageServerLoad } from './$types';
import { redirect, type Actions, fail, error } from '@sveltejs/kit';
import { z } from 'zod';
import { message, superValidate } from 'sveltekit-superforms';
import { dev } from '$app/environment';
import { zod } from 'sveltekit-superforms/adapters';

const schema = z.object({
	username: z.string().min(3).max(28),
});

//TODO custom error message
//TODO should say "Username must contain at least 3 character(s)" instead of "String must contain at least 3 character(s)"
export const load: PageServerLoad = async ({ locals: { supabase }, url }) => {
	const onboard = url.searchParams.get('onboard');

	const {
		data: { session },
	} = await supabase.auth.getSession();
	if (!session) {
		redirect(308, '/');
	}

	if (onboard) {
		const fetchUrl = dev
			? `http://${import.meta.env.VITE_API_URL}/user/onboarding`
			: `https://${import.meta.env.VITE_API_URL}/user/onboarding`;
		const token = session.access_token;
		const options = {
			method: 'PUT',
			headers: {
				Authorization: `Bearer ${token}`,
			},
		};
		fetch(fetchUrl, options);
	}

	const form = await superValidate(zod(schema));
	return { form };
};

export const actions: Actions = {
	default: async ({ locals: { supabase }, request }) => {
		const form = await superValidate(request, zod(schema));
		if (!form.valid) {
			return fail(400, { form });
		}

		const {
			data: { session },
		} = await supabase.auth.getSession();
		if (!session) {
			redirect(308, '/');
		}

		const fetchUrl = dev
			? `http://${import.meta.env.VITE_API_URL}/user/changename`
			: `https://${import.meta.env.VITE_API_URL}/user/changename`;
		const token = session.access_token;
		const options = {
			method: 'PUT',
			headers: {
				Authorization: `Bearer ${token}`,
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ username: form.data.username }),
		};
		const res = await fetch(fetchUrl, options);
		if (!res.ok) {
			error(500);
		}
		await supabase.auth.refreshSession(session);

		//TODO unique username validation on backend
		// return setError(form, 'username', 'Username already exists')

		return message(form, 'changed');
	},
};
