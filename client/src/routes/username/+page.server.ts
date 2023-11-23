// import { AuthApiError } from '@supabase/supabase-js';
import type { PageServerLoad } from './$types';
import { redirect, type Actions, fail, error } from '@sveltejs/kit';
import { z } from 'zod';
import { message, superValidate } from 'sveltekit-superforms/server';

const schema = z.object({
	username: z.string().min(3).max(28),
});

//TODO custom error message
//TODO should say "Username must contain at least 3 character(s)" instead of "String must contain at least 3 character(s)"
export const load: PageServerLoad = async ({ locals: { getSession }, url }) => {
	const onboard = url.searchParams.get('onboard');

	const session = await getSession();
	if (!session) {
		throw redirect(308, '/');
	}

	if (onboard) {
		const fetchUrl = `http://${import.meta.env.VITE_API_URL}/user/onboarding`;
		const token = session.access_token;
		const options = {
			method: 'PUT',
			headers: {
				Authorization: `Bearer ${token}`,
			},
		};
		fetch(fetchUrl, options);
	}

	const form = await superValidate(schema);
	return { form };
};

export const actions: Actions = {
	default: async ({ locals, request }) => {
		const form = await superValidate(request, schema);
		if (!form.valid) {
			return fail(400, { form });
		}

		const session = await locals.getSession();
		if (!session) {
			throw redirect(308, '/');
		}

		const fetchUrl = `http://${import.meta.env.VITE_API_URL}/user/changename`;
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
			throw error(500)
		}
		await locals.supabase.auth.refreshSession(session)

		//TODO unique username validation on backend
		// return setError(form, 'username', 'Username already exists')

		return message(form, 'changed');
	},
};
