import { getContext, setContext } from 'svelte';

function setService<T>(key: string, service: T): T {
	setContext(key, service);
	return service;
}

function getService<T>(key: string): () => T {
	return () => getContext(key) as T;
}

export function createService<T>(key: string): { get: () => T; set: (service: T) => T } {
	return {
		get: getService<T>(key),
		set: (service: T) => {
			setService(key, service);
			return service;
		},
	};
}
