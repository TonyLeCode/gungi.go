// import { redirect } from '@sveltejs/kit';
// import type { RequestHandler } from './$types';

// export const GET: RequestHandler = async (event) => {
// 	const {
// 		url,
// 		locals: { supabase },
// 	} = event;
// };

import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals: { getSession } }) => {
	// console.log("get session: ", await getSession())
	return {
		session: await getSession(),
	};
};
