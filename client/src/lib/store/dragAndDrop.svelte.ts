// I attempted to make drag and drop for mobile but there
// was no suitable support for dropping on touch devices

type dropOptions<T> = {
	mouseEnterEvent?: () => void;
	mouseLeaveEvent?: () => void;
	mouseEnterItem?: T;
};

export class Droppable<T> {
	hoverItem = $state<T | null>(null);
	isDragging = $state(false);
	constructor() {}

	addDroppable(node: HTMLElement, options?: dropOptions<T>) {
		function mouseLeave(this: Droppable<T>) {
			if (!this.isDragging) return;
			if (options?.mouseLeaveEvent) {
				options.mouseLeaveEvent();
			}
			if (options?.mouseEnterItem) {
				this.hoverItem = null;
			}
		}
		function mouseEnter(this: Droppable<T>) {
			if (!this.isDragging) return;
			if (options?.mouseEnterEvent) {
				options.mouseEnterEvent();
			}
			if (options?.mouseEnterItem) {
				this.hoverItem = options.mouseEnterItem;
			}
		}
		node.addEventListener('mouseleave', mouseLeave.bind(this));
		node.addEventListener('mouseenter', mouseEnter.bind(this));
	}
	startDragging() {
		this.isDragging = true;
	}

	stopDragging() {
		this.isDragging = false;
		this.hoverItem = null;
	}
}

const timeThreshold = 200;
let initialX = 0;
let initialY = 0;
let offsetX = 0;
let offsetY = 0;
let dragStart = true;
let clickTimeout: number | undefined;
let longPress = false;
let startTime: number | undefined;
let unsubMoveHandler: () => void;
let unsubReleaseHandler: () => void;

export type DraggableOptions<T> = {
	startEvent?: (hoverItem?: T) => void;
	dragEvent?: (hoverItem?: T) => void;
	dragStartEvent?: (hoverItem?: T) => void;
	releaseEvent?: (hoverItem?: T) => void;
	longReleaseEvent?: (hoverItem?: T) => void;
	shortReleaseEvent?: (hoverItem?: T) => void;
	dragReleaseEvent?: (hoverItem?: T) => void;
	// setDragItem?: T;
	droppable?: Droppable<T>;
	active?: boolean | ((hoverItem?: T) => boolean);
};

export function draggable<T>(node: HTMLElement, options: DraggableOptions<T> = {}) {
	const {
		startEvent,
		dragEvent,
		dragStartEvent,
		releaseEvent,
		longReleaseEvent,
		shortReleaseEvent,
		dragReleaseEvent,
		// setDragItem,
		droppable,
		active = true,
	} = options;

	if (typeof active === 'function') {
		if (!active()) return;
	} else if (active === false) return;

	function dragMoveHandler(node: HTMLElement) {
		return function (e: MouseEvent) {
			if (dragStart) {
				dragStart = false;
				if (droppable) {
					droppable.startDragging();
				}
				if (typeof dragStartEvent === 'function') {
					dragStartEvent(droppable?.hoverItem ?? undefined);
				}
			}
			if (typeof dragEvent === 'function') {
				dragEvent(droppable?.hoverItem ?? undefined);
			}
			let dx = 0;
			let dy = 0;

			if (clickTimeout !== undefined) {
				clearTimeout(clickTimeout);
				clickTimeout = undefined;
				startTime = undefined;
				longPress = false;
			}

			dx = e.clientX + offsetX - initialX;
			dy = e.clientY + offsetY - initialY;

			node.style.left = `${dx}px`;
			node.style.top = `${dy}px`;
			node.style.zIndex = '3';
		};
	}
	function dragReleaseHandler(node: HTMLElement) {
		return function () {
			if (startTime !== undefined) {
				const endTime = Date.now();
				const duration = endTime - startTime;
				if (duration < timeThreshold) {
					// short press
					if (typeof shortReleaseEvent === 'function') {
						shortReleaseEvent(droppable?.hoverItem ?? undefined);
					}
				} else if (longPress) {
					// long press
					if (typeof longReleaseEvent === 'function') {
						longReleaseEvent(droppable?.hoverItem ?? undefined);
					}
				}
			} else {
				// dragged
				if (typeof dragReleaseEvent === 'function') {
					dragReleaseEvent(droppable?.hoverItem ?? undefined);
				}
			}

			dragStart = true;
			clickTimeout = undefined;
			startTime = undefined;
			longPress = false;

			if (typeof releaseEvent === 'function') {
				releaseEvent(droppable?.hoverItem ?? undefined);
			}

			if (droppable) {
				droppable.stopDragging();
			}

			node.removeAttribute('style');
			initialX = 0;
			initialY = 0;
			offsetX = 0;
			offsetY = 0;
			unsubMoveHandler();
			unsubReleaseHandler();
		};
	}
	function dragStartHandler(e: MouseEvent) {
		if (e.target === null) return;

		if (typeof startEvent === 'function') {
			startEvent(droppable?.hoverItem ?? undefined);
		}

		startTime = Date.now();
		clickTimeout = window.setTimeout(() => {
			longPress = true;
		}, timeThreshold);

		const target = e.target as HTMLElement;
		if (e.button !== 0) return;
		initialX = e.clientX;
		initialY = e.clientY;
		offsetX = e.offsetX - target?.offsetWidth / 2;
		offsetY = e.offsetY - target?.offsetHeight / 2;

		const onDrag = dragMoveHandler(target);
		const onRelease = dragReleaseHandler(target);
		node.style.pointerEvents = 'none';
		document.addEventListener('mousemove', onDrag);
		document.addEventListener('mouseup', onRelease);
		unsubMoveHandler = () => {
			document.removeEventListener('mousemove', onDrag);
		};
		unsubReleaseHandler = () => {
			document.removeEventListener('mouseup', onRelease);
		};
	}

	node.addEventListener('mousedown', dragStartHandler);
	return {
		destroy() {
			node.removeEventListener('mousedown', dragStartHandler);
			// Note that destruction can happen while another drag is in progress
			if (unsubMoveHandler) unsubMoveHandler();
			if (unsubReleaseHandler) unsubReleaseHandler();
		},
	};
}