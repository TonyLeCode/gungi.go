<script lang="ts">
	import { onMount } from 'svelte';
	import CreateGameDialogue from './CreateGameDialogue.svelte';
	import RoomDialogue from './RoomDialogue.svelte';
	import RoomList from './RoomList.svelte';
	type Info = {
		name: string;
		description: string;
		type: string;
		color: string;
		rules: string;
	};

	let roomList: Info[] = [
		{
			name: 'Ornable',
			description:
				'DescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescription',
			type: 'correspondence',
			color: 'random',
			rules: 'default',
		},
		{
			name: 'Madahachi',
			description:
				'DescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescription',
			type: 'correspondence',
			color: 'random',
			rules: 'default',
		},
		{
			name: 'test',
			description:
				'DescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescription',
			type: 'live',
			color: 'random',
			rules: 'default',
		},
	];
	$: liveRoomList = roomList.filter((room) => {
		return room.type === 'live';
	});
	$: correspondenceRoomList = roomList.filter((room) => {
		return room.type === 'correspondence';
	});

	let showCreateGameDialogue = false;
	let showRoomDialogue = false;
	let roomDialogueInfo: Info = {
		name: 'test',
		description: 'DescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescription',
		type: 'correspondence',
		color: 'random',
		rules: 'default',
	};
	const url = (route: string) => `ws://${import.meta.env.VITE_API_URL}/${route}`;

	let text = 'not connected yet...';
	onMount(() => {
		const ws = new WebSocket(url('room'));

		ws.addEventListener('open', (event) => {
			text = 'connected!';
			setTimeout(() => {
				ws.send('hello!');
			}, 1750);
		});
		ws.addEventListener('message', (event) => {
			console.log('got message! ', event);
			// const data = JSON.parse(event.data);
			// console.log(data);
			text = event.data;
		});
	});
</script>

<main>
	<div>
		{text}
	</div>

	<div>
		<button
			on:click={() => {
				showCreateGameDialogue = true;
			}}
			class="button-primary">Create Game</button
		>
	</div>

	<h2>Live Games</h2>
	<RoomList bind:showRoomDialogue bind:roomDialogueInfo roomList={liveRoomList} />
	<h2>Correspondence Games</h2>
	<RoomList bind:showRoomDialogue bind:roomDialogueInfo roomList={correspondenceRoomList} />
	<CreateGameDialogue bind:showModal={showCreateGameDialogue} />
	<RoomDialogue bind:showModal={showRoomDialogue} info={roomDialogueInfo} />
</main>

<style lang="scss">
	main {
		max-width: 70rem;
		margin: auto;
		padding: 0.5rem;
	}
</style>
