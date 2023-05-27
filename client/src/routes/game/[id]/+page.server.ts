import { VITE_API_URL } from './../../../../.svelte-kit/ambient.d';
import axios from 'axios';

export async function load({params}){
  const url = `http://${import.meta.env.VITE_API_URL}/getgame/${params.id}`
  const data = await axios({
		method: 'get',
		url: url,
	}).then((res) => {
		return res.data;
	});
  return {
    data
  }
}