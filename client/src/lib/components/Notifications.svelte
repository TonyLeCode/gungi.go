<script lang="ts">
	import { AddNotification, notifications, type notificationType } from '$lib/store/notification';
  import { onMount } from 'svelte';
	import NotificationItem from './NotificationItem.svelte';
	import { nanoid } from 'nanoid';
  import { fly } from 'svelte/transition';

  onMount(() => {
    AddNotification({
			id: nanoid(),
      title: "Game Accepted",
			type: 'default',
			msg: 'Go to <a class="a-primary" href="/play/lol">game<a>',
		} as notificationType)
  })
</script>

<ul>
	{#each $notifications as notification (notification.id)}
		<li transition:fly|global="{{x: "16px", duration: 250}}"><NotificationItem notification={notification} /></li>
	{/each}
</ul>

<style lang="scss">
  ul{
    position: fixed;
    bottom: 2rem;
    right: 2rem;
    z-index: 10;
    // border: 1px red solid;
    display: flex;
    flex-direction: column;
    gap: .5rem;
  }
</style>
