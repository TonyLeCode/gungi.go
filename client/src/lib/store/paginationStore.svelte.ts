// Page starts at 1
export class PaginationStore {
	currentPage = $state(1);
	totalPages = $state<number>(1);
	hasNext = $derived(this.currentPage < this.totalPages);
	hasPrev = $derived(this.currentPage > 1);

	constructor(totalPages: number) {
		if (totalPages < 1) {
			throw new Error('totalPages must be greater than 0');
		}
		this.totalPages = totalPages;
	}

	setPage(page: number) {
		this.currentPage = page;
	}
  setTotalPages(totalPages: number) {
    this.totalPages = totalPages;
		if (this.currentPage > totalPages) {
			this.currentPage = totalPages;
		}
  }
	prev() {
		if (this.hasPrev) {
			this.currentPage -= 1;
		}
	}
	next() {
		if (this.hasNext) {
			this.currentPage += 1;
		}
	}
}

export function createPaginationStore(totalPages: number) {
  const paginationStore = new PaginationStore(totalPages);
	return paginationStore;
}
