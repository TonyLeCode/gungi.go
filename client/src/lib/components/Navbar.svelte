<script lang="ts">
	import type { Session } from '@supabase/supabase-js';
	import { Menu, SunMoon, X } from 'lucide-svelte';
	import { onMount } from 'svelte';

	let { session }: { session: Session | null } = $props();

	let isMenuOpen = $state(false);

	let theme = $state("light");

	function switchTheme() {
		if (theme === 'dark') {
			localStorage.setItem('theme', 'light');
			theme = 'light';
			document.documentElement.classList.remove('dark');
		} else {
			localStorage.setItem('theme', 'dark');
			theme = 'dark';
			document.documentElement.classList.add('dark');
		}
	}

	function resizeHandler() {
		const mediaQuerySize = 767;
		if (window.innerWidth >= mediaQuerySize) {
			isMenuOpen = false;
		}
	}

	function onClick(){
		isMenuOpen = false
	}

	onMount(() => {
		theme = localStorage.getItem('theme') ?? 'light';
		window.addEventListener('resize', resizeHandler);
		return () => {
			window.removeEventListener('resize', resizeHandler);
		};
	});
</script>

<nav class="navbar">
	<a onclick={onClick} class="brand" href="/"><img src="/gungi-logo.svg" alt="logo" /></a>
	<div class="nav-inner" class:open={isMenuOpen}>
		<ul class="nav-list" class:open={isMenuOpen}>
			{#if session}
				<a onclick={onClick} href="/overview">overview</a>
				<a onclick={onClick} href="/play">play</a>
				<!-- <a href="/rules">rules</a> -->
				<!-- TODO learning, puzzles, resources, library -->
			{/if}
			<!-- <a href="/rules">rules</a> -->
		</ul>
		<ul class="nav-options" class:open={isMenuOpen}>
			<button class={`theme-switcher ${theme}`} class:open={isMenuOpen} onclick={switchTheme}><SunMoon /></button>
			{#if session}
				<!-- TODO Dropdown for: profile, settings, friends, notifications  -->
				<div class="name">{session.user.user_metadata.username}</div>
				<a onclick={onClick} href="/logout">logout</a>
			{:else}
				<a onclick={onClick} class="a" href="/login">login</a>
				<a onclick={onClick} class="a" href="/register">register</a>
			{/if}
		</ul>
	</div>
	<button class="nav-menu" onclick={() => (isMenuOpen = !isMenuOpen)}>
		{#if isMenuOpen}
			<X color="rgb(var(--primary))" size="30px" />
		{:else}
			<Menu color="rgb(var(--primary))" size="30px" />
		{/if}
	</button>
</nav>

<style>
	nav {
		user-select: none;
		font-size: 0.85rem;
		@media (min-width: 767px) {
			font-size: 1rem;
		}
	}

	.navbar {
		position: sticky;
		top: 0;
		display: flex;
		max-width: 96rem;
		margin: auto;
		align-items: center;
		padding: 0.25rem 0.25rem;
		background-color: rgb(var(--bg));
		z-index: 4;
		/* max-width: 120rem; */
		/* margin: auto; */
		@media (min-width: 767px) {
			padding: 0.75rem 1rem;
		}
	}

	a {
		margin: 0;
		padding: 0;
		@media (min-width: 767px) {
			padding: 0 0.25rem;
			margin: 0 0.25rem;
		}
	}
	.brand {
		/* font-weight: 600; */
		/* width: 30px; */
		max-width: 30px;
		max-height: 30px;
		width: 100%;
		@media (min-width: 767px) {
			max-width: 45px;
			max-height: 45px;
		}
	}
	.name {
		display: none;
		color: rgb(var(--primary));
		margin-right: 1rem;
		position: relative;
		font-weight: 600;
		@media (min-width: 767px) {
			display: block;
		}
	}

	.nav-inner {
		display: flex;
		width: 100%;
	}

	.nav-inner.open {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: rgba(var(--bg), 0.96);
		z-index: 4;
		gap: 2rem;
		padding: 1rem;
		font-size: 1.2rem;
		padding: 1rem;
		margin-top: 38px;
		flex-direction: column;
	}
	.nav-list {
		display: none;
		gap: 0.5rem;
		@media (min-width: 767px) {
			display: flex;
			align-items: center;
		}
	}
	.nav-list.open {
		display: flex;
		flex-direction: column;
	}
	.nav-options {
		display: none;
		gap: 1rem;
		@media (min-width: 767px) {
			display: flex;
			align-items: center;
			margin-left: auto;
		}
	}

	.nav-options.open {
		display: flex;
		flex-direction: column;
	}

	.theme-switcher {
		size: 20px;
		display: flex;
		gap: 1rem;
		@media (min-width: 767px) {
			padding: 0 0.25rem;
		}
		&:hover {
			color: rgb(var(--primary));
		}
	}

	.theme-switcher.open.dark::before {
		content: 'light theme';
		display: block;
	}
	.theme-switcher.open.light::before {
		content: 'dark theme';
		display: block;
	}

	/* .nav-list.nav-open {
		display: flex;
		flex-direction: column;
		position: fixed;
		top: 0;
		left: 0;
		bottom: 0;
		right: 0;
		background-color: rgba(var(--bg), 0.96);
		color: rgb(var(--font));
		z-index: 4;
		gap: 1rem;
		font-size: 1.2rem;
		padding: 1rem;
		margin-top: 38px;
	} */

	.nav-menu {
		margin-left: auto;
		@media (min-width: 767px) {
			display: none;
		}
	}
</style>
