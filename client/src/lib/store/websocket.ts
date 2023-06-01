import { writable } from 'svelte/store';

type wsConnStateType = 'connecting' | 'connected' | 'reconnecting' | 'closed' | 'error';
type wsRouteType = 'overview' | 'game' | 'roomList';

export const ws = writable<WebSocket>();
export const wsConnState = writable<wsConnStateType>('closed');
export const wsRoute = writable<wsRouteType>();

export function websocketConnect(url: string, token: string) {
	const newWS = new WebSocket(url);

	wsConnState.set('connecting');

	newWS.addEventListener('open', () => {
		const msg = {
			type: 'auth',
			payload: `Bearer ${token}`,
		};
		// wsConnState.set('connected');
		newWS.send(JSON.stringify(msg));
	});

	newWS.addEventListener('message', (event) => {
		try {
			const data = JSON.parse(event.data);
			if (data.type == 'auth') {
				data.payload == '1' ? wsConnState.set('connected') : newWS.close();
			}
		} catch (err) {
			console.error('Error: ', err);
		}
	});

	newWS.addEventListener('error', (event) => {
		console.log(event);
		console.error('Error: ', event);
		wsConnState.set('error');
	});

	newWS.addEventListener('close', (event) => {
		console.log(event);
		wsConnState.set('closed');
	});

	ws.set(newWS);
}
