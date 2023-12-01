import type { Actions } from './$types';
import { AuthApiError } from '@supabase/supabase-js';
import type { PageServerLoad } from './$types';
import { fail, redirect, error } from '@sveltejs/kit';
import { z } from 'zod';
import { superValidate, message } from 'sveltekit-superforms/server';
import { dev } from '$app/environment';

const schema = z.object({
	email: z.string(),
	password: z.string(),
});

export const load: PageServerLoad = async ({ locals: { getSession } }) => {
	const session = await getSession();
	if (session) {
		throw redirect(308, '/overview');
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

		const supabaseResponse = await locals.supabase.auth.signInWithPassword({
			email: form.data.email,
			password: form.data.password,
		});

		if (supabaseResponse.error) {
			console.log(supabaseResponse.error);
			if (supabaseResponse.error instanceof AuthApiError && supabaseResponse.error.status === 400) {
				if (supabaseResponse.error.message === "Email not confirmed") {
					return message(form, 'Email not verified');
				}
				return message(form, 'Invalid login info');
			}
			throw error(500, {
				message: 'Server error. Try again later.',
			});
		}
		const fetchUrl = dev ? `http://${import.meta.env.VITE_API_URL}/user/onboarding` : `https://${import.meta.env.VITE_API_URL}/user/onboarding`;
		const token = supabaseResponse.data.session.access_token
		const options = {
			headers: {
				Authorization: `Bearer ${token}`,
			},
		};
		const res = await fetch(fetchUrl, options);
		if (res.ok) {
			const hasOnboarded = await res.json();
			if (!hasOnboarded) throw redirect(308, '/username?onboard=true')
		}

		// throw redirect(303, '/overview');
	},
} satisfies Actions;
