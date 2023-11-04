// import { readable } from 'svelte/store';

// export const x = readable({ x: 0, y: 0 }, (set) => {
// 	function move(event: MouseEvent) {
// 		set({
// 			x: event.clientX,
// 			y: event.clientY,
// 		});
// 	}
// 	document.body.addEventListener('mousemove', move);

// 	return () => {
// 		document.body.removeEventListener('mousemove', move);
// 	};
// });

// export function act(node: HTMLElement) {
// 	let unsub: any;
// 	function start() {
// 		unsub = x.subscribe((value) => {
// 			console.log(value);
// 		});
// 	}
// 	function up() {
// 		if (unsub) {
// 			unsub();
// 		}
// 	}
// 	node.addEventListener('mousedown', start);
// 	document.addEventListener('mouseup', up);
// 	return {
// 		destroy() {
// 			node.removeEventListener('mousedown', start);
// 			document.removeEventListener('mouseup', up);
// 		},
// 	};
// }
