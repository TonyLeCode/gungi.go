<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import CreateGameDialogue from './CreateGameDialogue.svelte';
	import RoomDialogue from './RoomDialogue.svelte';
	import RoomList from './RoomList.svelte';
	import { AddNotification, type notificationType } from '$lib/store/notification';
	import {nanoid} from 'nanoid'
	import {ws, wsConnState} from '$lib/store/websocket'

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

	let roomList: Info[] = [];
	$: roomList = roomList.sort((a, b) => (a.host === username ? -1 : 1));
	let showLive = true;
	let showCorrespondence = true;
	$: liveRoomList = roomList?.filter((room) => room.type === 'live');
	$: correspondenceRoomList = roomList?.filter((room) => room.type === 'correspondence');

	let showCreateGameDialogue = false;
	let showRoomDialogue = false;
	let roomDialogueInfo: Info;

	function handleRoomListMsg(event: MessageEvent<any>){
		console.log(event);
			try {
				const data = JSON.parse(event.data);
				switch (data.type) {
					case 'roomList':
						roomList = JSON.parse(data.payload);
						break;
					case 'accepted':
						console.log(data.payload);
						AddNotification({
							id: nanoid(),
							title: 'Game Accepted',
							type: 'default',
							msg: `Go to <a class="a-primary" href="/game/${data.payload}">game<a>`,
						} as notificationType);
						break;
				}
			} catch (err) {
				console.log(event?.data);
				console.error('Error: ', err);
			}
	}

	function accept(roomid: string) {
		const msg = {
			type: 'accept',
			payload: roomid,
		};
		$ws.send(JSON.stringify(msg));
	}

	onMount(() => {
		$ws.addEventListener('message', handleRoomListMsg);

		return () => {
			$ws.removeEventListener('message', handleRoomListMsg)
		};
	});

</script>

<main>
	{#if $wsConnState === 'connecting'}
		<p class="status-msg fly-up-fade">Loading...</p>
	{:else if $wsConnState === 'connected'}
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
					class="button-primary">Create Game</button
				>
			</div>
		</section>
		{#if showLive}
			<RoomList
				bind:showRoomDialogue
				bind:roomDialogueInfo
				ws={$ws}
				roomList={liveRoomList}
				heading="Live Games"
				{username}
			/>
		{/if}
		{#if showCorrespondence}
			<RoomList
				bind:showRoomDialogue
				bind:roomDialogueInfo
				ws={$ws}
				roomList={correspondenceRoomList}
				heading="Correspondence Games"
				{username}
			/>
		{/if}
		<CreateGameDialogue bind:showModal={showCreateGameDialogue} host={username} ws={$ws} />
		<RoomDialogue bind:showModal={showRoomDialogue} info={roomDialogueInfo} {accept} />
	{:else if $wsConnState === 'error'}
		<p class="status-msg fly-up-fade">Something went wrong, please refresh or try again later</p>
	{:else if $wsConnState === 'closed'}
		<p class="status-msg fly-up-fade">Not connected, please refresh or try again later</p>
	{/if}
</main>

<style lang="scss">
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
	.status-msg {
		margin-top: 3rem;
		text-align: center;
	}
	.filter {
		justify-content: center;
		align-items: center;
		display: flex;
		gap: 1rem;
		margin-left: auto;
	}
</style>
