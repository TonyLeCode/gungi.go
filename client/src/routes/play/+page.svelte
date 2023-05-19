<script lang="ts">
	import { onMount } from 'svelte';
	import CreateGameDialogue from './CreateGameDialogue.svelte';
	import RoomDialogue from './RoomDialogue.svelte';
	import RoomList from './RoomList.svelte';
  import { fly } from 'svelte/transition';

	export let data;

	$: username = data.session?.user.user_metadata.username;

	type Info = {
		roomid: string;
		host: string;
		description: string;
		type: string;
		color: string;
		rules: string;
	};

	// let roomList: Info[] = [
	// 	{
	// 		roomid: '',
	// 		host: 'Ornable',
	// 		description:
	// 			'DescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescription',
	// 		type: 'correspondence',
	// 		color: 'random',
	// 		rules: 'default',
	// 	},
	// 	{
	// 		roomid: '',
	// 		host: 'Madahachi',
	// 		description:
	// 			'DescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescription',
	// 		type: 'correspondence',
	// 		color: 'random',
	// 		rules: 'default',
	// 	},
	// 	{
	// 		roomid: '',
	// 		host: 'test',
	// 		description:
	// 			'DescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescription',
	// 		type: 'live',
	// 		color: 'random',
	// 		rules: 'default',
	// 	},
	// ];
	let roomList: Info[] = [];
	let showLive = true;
	let showCorrespondence = true;

	$: liveRoomList = roomList?.filter((room) => {
		return room.type === 'live';
	});
	$: correspondenceRoomList = roomList?.filter((room) => {
		return room.type === 'correspondence';
	});

	let showCreateGameDialogue = false;
	let showRoomDialogue = false;
	let roomDialogueInfo: Info;

	const url = (route: string) => `ws://${import.meta.env.VITE_API_URL}/${route}`;

	let text = 'not connected yet...';
	let ws: WebSocket;
	onMount(() => {
		ws = new WebSocket(url('room'));

		ws.addEventListener('open', (event) => {
			text = 'connected!';
			// setTimeout(() => {
			// 	ws.send(JSON.stringify(payload));
			// }, 1750);
		});
		ws.addEventListener('message', (event) => {
			// console.log('got message! ', event);
			const data = JSON.parse(event.data);
			if (data?.type === 'roomList') {
				roomList = JSON.parse(data.payload);
			}
			console.log('got message! ', data);
			// const data = JSON.parse(event.data);
			// console.log(data);
		});
	});
</script>

<main>
	<div>
		{text}
	</div>
	<section class="options">
		<div class="filter">
			<label>
				<input bind:checked={showLive} type="checkbox" />
				Live
			</label>
			<label>
				<input bind:checked={showCorrespondence} type="checkbox" />
				Correspondence
			</label>
			<button
				on:click={() => {
					showCreateGameDialogue = true;
				}}
				class="create-game button-primary">Create Game</button
			>
		</div>
	</section>

	{#if showLive}
		<h2 class="fly-up">Live Games</h2>
		<RoomList bind:showRoomDialogue bind:roomDialogueInfo roomList={liveRoomList} />
	{/if}
	{#if showCorrespondence}
		<h2 class="fly-up">Correspondence Games</h2>
		<RoomList bind:showRoomDialogue bind:roomDialogueInfo roomList={correspondenceRoomList} />
	{/if}
	<CreateGameDialogue bind:showModal={showCreateGameDialogue} host={username} {ws} />
	<RoomDialogue bind:showModal={showRoomDialogue} info={roomDialogueInfo} />
</main>

<style lang="scss">
	h2 {
		margin-bottom: 0.5rem;
	}
	main {
		max-width: 70rem;
		margin: auto;
		padding: 0.5rem;
		margin-top: 2rem;
	}
	.options {
		display: flex;
		margin: 0 2rem;
	}
	.create-game {
		// margin-left: auto;
	}
	.filter {
		justify-content: center;
		align-items: center;
		display: flex;
		gap: 1rem;
		margin-left: auto;
	}
</style>
