import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url, locals: { getSession } }) => {
	const session = await getSession();
	// console.log('session', session)
	console.log('url', url.origin);

	if (session) {
		// throw redirect(303, '/');
	}

	return { url: url.origin };
};
