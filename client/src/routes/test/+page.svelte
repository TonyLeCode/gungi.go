<script>
	// import {supabase} from '$lib/supabaseClient'

	export let data;
	let { countries } = data;
	// $: ({countries} = data);
	// const insertRow = async () => {
	//   const {data, error} = await supabase
	//     .from('countries')
	//     .insert({name:'China'})
	//   if (error){
	//     console.log(error.message)
	//   } else {
	//     console.log(data)
	//   }
	// }
	// insertRow()
</script>

<ul>
	{#each countries as country}
		<li>{country.name}</li>
	{/each}
</ul>
