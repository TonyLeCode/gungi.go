<script lang="ts">
	import type { Session } from '@supabase/supabase-js';

	export let session: Session | null;
	$: session = session;
	// $: console.log("session!!", session)
</script>

<nav class="navbar">
	<a class="brand" href="/">Gungi</a>
	<ul class="nav-list">
		{#if session}
			<a href="/overview">overview</a>
			<a href="/play">play</a>
			<!-- <a href="/games">games</a> -->
			<!-- TODO learning, puzzles, resources, library -->
		{/if}
		<a href="/rules">rules</a>
	</ul>
	<ul class="nav-account">
		{#if session}
			<!-- TODO settings, friends, notifications  -->
			<span class="name">{session.user.user_metadata.username}</span>
			<a href="/logout">logout</a>
		{:else}
			<a class="a" href='login'>login</a>
			<a class="a" href='register'>register</a>
		{/if}
	</ul>
</nav>

<style>
	nav {
		user-select: none;
	}
	a {
		padding: 0 0.5rem;
	}
	.brand {
		font-weight: 600;
		/* font-size: 1.25rem; */
	}
	.name {
		color: rgb(var(--primary));
		margin-right: 1rem;
		position: relative;
		font-weight: 600;
	}
	/* .name::before{
		content: '';
		width: 15px;
		height: 15px;
		border-radius: 50%;
		background-color: rgb(var(--primary));
		display:block;
		position:absolute;
		left: -1.25rem;
		margin: auto;
		top: 0;
		bottom: 0;
	} */

	.navbar {
		display: flex;
		/* justify-content: space-between; */
		/* background-color:red; */
		/* max-width: 1000px; */
		padding: 0.75rem 12rem;
		max-width: 120rem;
		margin: auto;
	}
	.nav-list {
		display: flex;
		gap: 0.5rem;
		margin-left: 2rem;
	}
	.nav-account {
		display: flex;
		gap: 1rem;
		margin-left: auto;
	}
</style>
