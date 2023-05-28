export interface dragAndDropOptions {
	startEvent?: (items?: dragAndDropItems) => void;
	dragEvent?: (items?: dragAndDropItems) => void;
	releaseEvent?: (items?: dragAndDropItems) => void;
	setDragItem?: unknown;
	active?: boolean;
}
export interface dropOptions {
	mouseEnterEvent?: (items?: dragAndDropItems) => void;
	mouseLeaveEvent?: (items?: dragAndDropItems) => void;
	mouseEnterItem?: unknown;
}

export interface dragAndDropItems {
	dragItem?: unknown;
	hoverItem?: unknown;
}

export type dragAndDropFunction = (node: HTMLElement, options?: dragAndDropOptions) => {
  destroy: () => void;
};

export type dropFunction = (node: HTMLElement, options?: dropOptions) => {
  destroy: () => void;
};

let initX = 0;
let initY = 0;
let offsetX = 0;
let offsetY = 0;
let dragElement: HTMLElement | null;
let dragItem: unknown;
let hoverItem: unknown;

export function drop(node: HTMLElement, options = {} as dropOptions) {
	const { mouseEnterEvent, mouseLeaveEvent, mouseEnterItem } = options;
	function mouseLeave() {
		if (dragElement && mouseLeaveEvent && typeof mouseLeaveEvent === 'function') {
			const items = {
				dragItem: dragItem,
				hoverItem: hoverItem
			} as dragAndDropItems
			mouseLeaveEvent(items);
		}
		if (dragElement) {
			hoverItem = null;
		}
	}
	function mouseEnter() {
		if (dragElement && mouseEnterItem != null) {
			hoverItem = mouseEnterItem;
		}

		if (dragElement && mouseEnterEvent && typeof mouseEnterEvent === 'function') {
			const items = {
				dragItem: dragItem,
				hoverItem: hoverItem
			} as dragAndDropItems
			mouseEnterEvent(items);
		}
	}
	node.addEventListener('mouseleave', mouseLeave);
	node.addEventListener('mouseenter', mouseEnter);

	return {
		destroy() {
			node.removeEventListener('mouseleave', mouseLeave);
			node.removeEventListener('mouseenter', mouseEnter);
		},
	};
}

export function dragAndDrop(node: HTMLElement, options = {} as dragAndDropOptions) {
	const { startEvent, dragEvent, releaseEvent, setDragItem, active = true } = options;

	function releaseMouse() {
		if (dragElement) {
			dragElement.removeAttribute("style")
			document.removeEventListener('mousemove', dragMouse);
			document.removeEventListener('mouseup', releaseMouse);

			if (releaseEvent && typeof releaseEvent === 'function') {
				const items = {
					dragItem: dragItem,
					hoverItem: hoverItem
				} as dragAndDropItems
				releaseEvent(items);
			}
			dragElement = null;
			hoverItem = null;
			if(setDragItem) {
				dragItem = null
			}
		}
	}

	function dragMouse(e: MouseEvent) {
		if (dragElement) {
			const dx = e.clientX + offsetX - initX;
			const dy = e.clientY + offsetY - initY;
			dragElement.style.left = `${dx}px`;
			dragElement.style.top = `${dy}px`;

			if (dragEvent && typeof dragEvent === 'function') {
				const items = {
					dragItem: dragItem,
					hoverItem: hoverItem
				} as dragAndDropItems
				dragEvent(items);
			}
		}
	}

	function startDragMouse(e: MouseEvent) {
		initX = e.clientX;
		initY = e.clientY;
		const target = e.target as HTMLElement;
		offsetX = e.offsetX - target?.offsetWidth / 2;
		offsetY = e.offsetY - target?.offsetHeight / 2;
		dragElement = target;
		dragElement.style.pointerEvents = 'none';
		dragElement.style.zIndex = '3';
		document.addEventListener('mousemove', dragMouse);
		document.addEventListener('mouseup', releaseMouse);

		if(setDragItem) {
			dragItem = setDragItem
		}

		if (startEvent && typeof startEvent === 'function') {
			const items = {
				dragItem: dragItem,
				hoverItem: hoverItem
			} as dragAndDropItems
			startEvent(items);
		}
	}

	active && node.addEventListener('mousedown', startDragMouse);
	return {
		destroy() {
			node.removeEventListener('mousedown', startDragMouse);
			document.removeEventListener('mousemove', dragMouse);
			document.removeEventListener('mouseup', releaseMouse);
		},
	};
}
