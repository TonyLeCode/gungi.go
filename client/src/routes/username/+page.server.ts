// import { AuthApiError } from '@supabase/supabase-js';
import type { PageServerLoad } from './$types';
import { redirect, type Actions, fail, error } from '@sveltejs/kit';
import { z } from 'zod';
import { setError, superValidate } from 'sveltekit-superforms';
import { dev } from '$app/environment';
import { zod } from 'sveltekit-superforms/adapters';

const schema = z.object({
	username: z
		.string()
		.min(3, { message: 'Username must contain at least 3 character(s)' })
		.max(28, { message: 'Username must contain less than 28 character(s)' })
		.trim()
		.regex(/^[a-zA-Z0-9 _-]+$/, {
			message: 'Username must only contain letters, numbers, spaces, hyphens, and underscores',
		}),
});

export const load: PageServerLoad = async ({ locals: { supabase, username } }) => {
	const {
		data: { session },
	} = await supabase.auth.getSession();
	if (!session) {
		redirect(308, '/');
	}

	if (username) {
		redirect(308, '/overview');
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
			? `http://${import.meta.env.VITE_API_URL}/username`
			: `https://${import.meta.env.VITE_API_URL}/username`;
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
			if (res.status === 409) {
				return setError(form, 'username', 'Username already exists');
			}
			error(500);
		}

		redirect(303, '/overview');
	},
};
