import type { Actions } from './$types';
import { AuthApiError } from '@supabase/supabase-js';
import type { PageServerLoad } from './$types';
import { fail, redirect, error } from '@sveltejs/kit';
import { z } from 'zod';
import { superValidate, message } from 'sveltekit-superforms/server';

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
				// return fail(400, {
				// 	error: "Invalid Login Info"
				// })
				// return fail(400, { form });
				return message(form, 'Invalid login info');
			}
			// return fail(500, {
			// 	error: 'Server error. Try again later.',
			// });
			throw error(500, {
				message: 'Server error. Try again later.',
			});
		} else {
			console.log('logged in', supabaseResponse.data);
			throw redirect(303, '/overview');
		}
	},
} satisfies Actions;
