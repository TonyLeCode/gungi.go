import { AuthApiError } from '@supabase/supabase-js';
import type { PageServerLoad } from './$types';
import { redirect, type Actions, fail } from '@sveltejs/kit';
import { z } from 'zod';
import { superValidate } from 'sveltekit-superforms/server';

const schema = z.object({
	username: z.string().min(3).max(24),
	email: z.string().email(),
	password: z.string().min(6).max(64),
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

		//TODO unique username validation on backend
		// return setError(form, 'username', 'Username already exists')

		const { data, error } = await locals.supabase.auth.signUp({
			email: form.data.email,
			password: form.data.password,
			options: {
				data: {
					username: form.data.username,
				},
			},
		});

		if (error) {
			console.log(error);
			if (error instanceof AuthApiError && error.status === 400) {
				return fail(400, {
					error: 'Invalid Registration',
				});
			}
			return fail(500, {
				error: 'Server error. Try again later.',
			});
		} else {
			console.log('registered: ', data);
		}

		return { form };
	},
};
