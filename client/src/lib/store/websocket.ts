import { readable } from 'svelte/store';
import { notifications, type notificationType } from './notification';
import { nanoid } from 'nanoid';
import { browser } from '$app/environment';

type wsConnStateType = 'connecting' | 'connected' | 'reconnecting' | 'closed' | 'error';

interface msgType {
	type: string;
	payload?: unknown;
}

function createWsStore() {
	if (!browser) return;
	const newSocket = new WebSocket(`${import.meta.env.VITE_API_URL}/ws`);
	const { subscribe } = readable<wsConnStateType>('closed', (set) => {
		newSocket.addEventListener('open', () => {
			set('connecting');
		});
		newSocket.addEventListener('error', () => {
			set('error');
		});
		newSocket.addEventListener('close', () => {
			set('closed');
			notifications?.add({
				id: nanoid(),
				title: 'Disconnected',
				type: 'default',
				msg: 'Please refresh or come back later',
			} as notificationType);
		});
		newSocket.addEventListener('message', (event) => {
			try {
				const data = JSON.parse(event.data);
				if (data.type == 'auth') {
					data.payload == 'success' ? set('connected') : newSocket.close();
				}
			} catch (err) {
				console.error('Error: ', err);
			}
		});
	});

	function send(msg: msgType) {
		newSocket.send(JSON.stringify(msg));
	}

	function authenticate(token: string) {
		const msg = {
			type: 'auth',
			payload: `Bearer ${token}`,
		};
		send(msg);
	}

	function close() {
		newSocket.close();
	}

	function addMsgListener(fn: (event?: MessageEvent) => void) {
		newSocket.addEventListener('message', fn);
		return () => {
			newSocket.removeEventListener('message', fn);
		};
	}

	console.log('create ws store');
	return {
		subscribe,
		send,
		authenticate,
		close,
		addMsgListener,
	};
}

export const ws = createWsStore();
