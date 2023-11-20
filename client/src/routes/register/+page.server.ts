import { AuthApiError } from '@supabase/supabase-js';
import type { PageServerLoad } from './$types';
import { redirect, type Actions, fail } from '@sveltejs/kit';
import { z } from 'zod';
import { message, superValidate } from 'sveltekit-superforms/server';
import { uniqueNamesGenerator, animals, colors, adjectives, NumberDictionary, type Config } from 'unique-names-generator';

const schema = z.object({
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

		const numberDictionary = NumberDictionary.generate({ length: 4 });
		const nameGenConfig: Config = {
			dictionaries: [adjectives, colors, animals, numberDictionary],
			separator: '',
			length: 4,
			style: 'capital',
		}
		const randomName = uniqueNamesGenerator(nameGenConfig)

		//TODO unique username validation on backend
		// return setError(form, 'username', 'Username already exists')

		const { data, error } = await locals.supabase.auth.signUp({
			email: form.data.email,
			password: form.data.password,
			options: {
				data: {
					// username: form.data.username,
					username: randomName,
				},
			},
		});

		if (error) {
			console.log(error);
			if (error instanceof AuthApiError && error.status === 400) {
				//TODO handle error here and on client side with messages
				return fail(400, {
					error: 'Invalid Registration',
				});
			}
			return fail(500, {
				error: 'Server error. Try again later.',
			});
		} else {
			console.log('registered: ', data);
			return message(form, "success")
		}
	},
};
