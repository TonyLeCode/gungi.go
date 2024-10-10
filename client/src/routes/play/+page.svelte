<script lang="ts">
	import { onMount } from 'svelte';
	import CreateGameDialog from './CreateGameDialog.svelte';
	import RoomDialog from './RoomDialog.svelte';
	import RoomList from './RoomList.svelte';
	import { getNotificationStore, type notificationType } from '$lib/store/notificationStore.svelte';
	import { nanoid } from 'nanoid';
	import { getWebsocketStore } from '$lib/store/websocketStore.svelte';
	import Modal from '$lib/components/Modal.svelte';

	// TODO make page responsive
	let { data } = $props();
	let notificationStore = getNotificationStore();
	let websocketStore = getWebsocketStore();

	let username = $derived(data.session?.user.user_metadata.username);

	type Info = {
		id: string;
		host: string;
		description: string;
		type: string;
		color: string;
		rules: string;
	};

	let roomList = $state<Info[]>([]);
	let sortedList = $derived(roomList.sort((a, _) => (a.host === username ? -1 : 1)));
	let showCorrespondence = $state(true);
	let correspondenceRoomList = $derived(sortedList.filter((room) => room.type === 'correspondence'));

	// let showCreateGameDialogue = $state(false);
	// let showRoomDialogue = $state(false);
	let createGameDialog = $state<ReturnType<typeof Modal>>();
	let roomDialog = $state<ReturnType<typeof Modal>>();
	let roomDialogueInfo = $state<Info>({ id: '', host: '', description: '', type: '', color: '', rules: '' });

	function handleRoomListMsg(event?: MessageEvent) {
		try {
			const data = JSON.parse(event?.data);
			switch (data.type) {
				case 'roomList':
					roomList = data.payload ?? [];
					break;
				case 'roomAccepted':
					notificationStore.add({
						id: nanoid(),
						title: 'Game Accepted',
						type: 'success',
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
			type: 'acceptPlayRoom',
			payload: roomid,
		};
		websocketStore.send(msg);
	}

	onMount(() => {
		$effect(() => {
			if (websocketStore.state === 'connected') {
				websocketStore.addMsgListener(handleRoomListMsg);
				const msg = {
					type: 'joinPlay',
				};

				websocketStore.send(msg);
			}
		});
		return () => {
			websocketStore.removeMsgListener(handleRoomListMsg);
			const msg = {
				type: 'leavePlay',
			};
			websocketStore.send(msg);
		};
	});
</script>

<svelte:head>
	<title>Room List | White Monarch Server</title>
</svelte:head>

<main>
	{#if websocketStore.state === 'connecting'}
		<p class="status-msg fly-up-fade">Loading...</p>
	{:else if websocketStore.state === 'connected'}
		<section class="options">
			<div class="filter">
				<!-- <label>
					<input bind:checked={showLive} type="checkbox" />
					Live
				</label>
				<label>
					<input bind:checked={showCorrespondence} type="checkbox" />
					Correspondence
				</label> -->
				<button
					onclick={() => {
						createGameDialog?.open();
					}}
					disabled={username == null}
					class="button-primary">Create Game</button
				>
			</div>
		</section>
		<!-- {#if showLive}
			<RoomList
				bind:showRoomDialogue
				bind:roomDialogueInfo
				roomList={liveRoomList}
				heading="Live Games"
				{username}
				{accept}
			/>
		{/if} -->
		{#if showCorrespondence}
			<RoomList
				bind:roomDialogueInfo
				roomList={correspondenceRoomList}
				heading="Correspondence Games"
				{username}
				{accept}
				{roomDialog}
			/>
		{/if}
		<CreateGameDialog bind:createGameDialog />
		<RoomDialog bind:roomDialog info={roomDialogueInfo} {accept} />
	{:else if websocketStore.state === 'error'}
		<p class="status-msg fly-up-fade">Something went wrong, please refresh or try again later</p>
	{:else if websocketStore.state === 'closed'}
		<p class="status-msg fly-up-fade">Not connected, please refresh or try again later</p>
	{/if}
</main>

<style lang="scss">
	main {
		max-width: 70rem;
		margin: auto;
		padding: 0.5rem;
		margin-top: 2rem;
		@media (min-width: 767px) {
			
		}
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
