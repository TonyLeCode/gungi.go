<script lang="ts">
	import { onMount } from 'svelte';

	const url = (route: string) => 'ws://' + '127.0.0.1:8080/' + route;

  let text = 'not connected yet...';
	onMount(() => {
		const ws = new WebSocket(url('room'));


		ws.addEventListener('open', (event) => {
      text = 'connected!'
			setTimeout(() => {
				ws.send('hello!');
			}, 1750);
		});
		ws.addEventListener('message', (event) => {
			console.log('got message! ', event);
			const data = JSON.parse(event.data);
			console.log(data);
			text = data.data;
		});
	});
</script>

<main>
	<div>
		{text}
	</div>
</main>

<style>
	main {
		max-width: 70rem;
		margin: auto;
    min-height: 60vh;
    display: flex;
    justify-content: center;
    align-items: center;
	}
	div {
		text-align: center;
		font-size: 4rem;
	}
</style>
