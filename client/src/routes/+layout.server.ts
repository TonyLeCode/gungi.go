import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals: { safeGetSession }, cookies }) => {
	const { session, username } = await safeGetSession();
	return {
		session,
		username,
		cookies: cookies.getAll(),
	};
};
