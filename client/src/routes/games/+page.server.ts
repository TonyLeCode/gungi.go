// import { getServerSession } from '@supabase/auth-helpers-sveltekit'

import type { PageServerLoad } from './$types';
import axios from 'axios';

export const load: PageServerLoad = async ({ locals: { getSession } }) => {
	const session = await getSession();
	const token = session?.access_token;
	console.log('token', session?.access_token);
	const url = 'http://localhost:5080/getongoinggamelist';
	// const token = session.
	const data = await axios({
		method: 'get',
		url: url,
		headers: {
			Authorization: `Bearer ${token}`,
		},
	}).then((res) => {
		return res.data;
	});
	// const data: string[] =[]
	// const { data } = await supabase.from("countries").select();
	// return {
	//   countries: data ?? [],
	// };
	return {
		data: data ?? [],
	};
};
