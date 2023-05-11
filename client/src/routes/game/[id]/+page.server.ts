import axios from 'axios';

export async function load({params}){
  const url = `http://localhost:5080/getgame/${params.id}`
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