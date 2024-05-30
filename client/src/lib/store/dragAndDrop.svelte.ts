// import { getContext, setContext } from 'svelte';

type dropOptions<T> = {
	mouseEnterEvent?: () => void;
	mouseLeaveEvent?: () => void;
	mouseEnterItem?: T;
};

export class Droppable<T> {
	hoverItem = $state<T | null>(null);
  isDragging = $state(false);
	constructor() {
  }

	addDroppable( node: HTMLElement, options?: dropOptions<T>) {
		function mouseLeave(this: Droppable<T>) {
      if (!this.isDragging) return;
      if (options?.mouseLeaveEvent){
        options.mouseLeaveEvent();
      }
      if (options?.mouseEnterItem){
        this.hoverItem = null;
      }
    }
		function mouseEnter(this: Droppable<T>) {
      if (!this.isDragging) return;
      if (options?.mouseEnterEvent){
        options.mouseEnterEvent();
      }
      if (options?.mouseEnterItem){
        this.hoverItem = options.mouseEnterItem;
      }
		}
		node.addEventListener('mouseleave', mouseLeave.bind(this));
		node.addEventListener('mouseenter', mouseEnter.bind(this));
	}
}

let initialX = 0;
let initialY = 0;
let offsetX = 0;
let offsetY = 0;
let unsubMoveHandler: () => void;
let unsubReleaseHandler: () => void;

type draggableOptions = {
  startEvent?: () => void;
  dragEvent?: () => void;
  releaseEvent?: () => void;
  setDragItem?: unknown;
  active?: boolean | (() => boolean);
};

export function draggable(node: HTMLElement, options: draggableOptions = {}) {
  const { startEvent, dragEvent, releaseEvent, setDragItem, active = true } = options;

	function dragMoveHandler(node: HTMLElement) {
		return function (e: MouseEvent | TouchEvent) {
			let dx = 0;
			let dy = 0;
			if (e instanceof TouchEvent) {
				e.preventDefault();
				dx = e.touches[0].clientX + offsetX - initialX;
				dy = e.touches[0].clientY + offsetY - initialY;
			} else if (e instanceof MouseEvent) {
				dx = e.clientX + offsetX - initialX;
				dy = e.clientY + offsetY - initialY;
			}
			node.style.left = `${dx}px`;
			node.style.top = `${dy}px`;
			node.style.zIndex = '3';
		};
	}
	function dragReleaseHandler(node: HTMLElement) {
		return function () {
			console.log('release');
			node.removeAttribute('style');
			initialX = 0;
			initialY = 0;
			offsetX = 0;
			offsetY = 0;
			unsubMoveHandler();
			unsubReleaseHandler();
		};
	}
	function dragStartHandler(e: MouseEvent | TouchEvent) {
    if (typeof active === 'function') {
      if (!active()) return;
    } else if (active === false) return;
		if (e.target === null) return;
    
		const target = e.target as HTMLElement;
		if (e instanceof TouchEvent) {
			initialX = e.touches[0].clientX;
			initialY = e.touches[0].clientY;
			const rect = target.getBoundingClientRect();
			offsetX = e.targetTouches[0].clientX - rect.x - target?.offsetWidth / 2;
			offsetY = e.targetTouches[0].clientY - rect.y - target?.offsetHeight / 2;
		} else if (e instanceof MouseEvent) {
			if (e.button !== 0) return;
			initialX = e.clientX;
			initialY = e.clientY;
			offsetX = e.offsetX - target?.offsetWidth / 2;
			offsetY = e.offsetY - target?.offsetHeight / 2;
		}
		const onDrag = dragMoveHandler(target);
		const onRelease = dragReleaseHandler(target);
		target.addEventListener('mousemove', onDrag);
		target.addEventListener('mouseup', onRelease);
		target.addEventListener('touchmove', onDrag);
		target.addEventListener('touchend', onRelease);
		unsubMoveHandler = () => {
			target.removeEventListener('mousemove', onDrag);
			target.removeEventListener('touchmove', onDrag);
		};
		unsubReleaseHandler = () => {
			target.removeEventListener('mouseup', onRelease);
			target.removeEventListener('touchend', onRelease);
		};
	}

	node.addEventListener('mousedown', dragStartHandler);
	node.addEventListener('touchstart', dragStartHandler);
}

// export function createDrag() {
// 	const store = new Drag();
// 	setContext('drag', store);
// 	return store;
// }

// export function getDrag() {
// 	return getContext('drag');
// }
