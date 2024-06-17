const baseUrl = 'https://ro-viz.c1.is/rovis/';

function buildSearchParams(searchParams: URLSearchParams, params?: Record<string, string | number>): void {
	if (!params) {
		return;
	}

	Object.entries(params).forEach(([key, val]) => {
		searchParams.set(key, val.toString());
	});
}

async function post<T>(path: string, body: object | null, params?: Record<string, string | number>): Promise<T> {
	const url = new URL(`${baseUrl}${path}`);
	buildSearchParams(url.searchParams, params);
	const res = await window.fetch(url, {
		method: 'POST',
		headers: new Headers({
			"Content-Type": "application/json",
		}),
		body: JSON.stringify(body),
	});

	return await res.json();
}

async function get<T>(path: string, params?: Record<string, string | number>): Promise<T> {
	const url = new URL(`${baseUrl}${path}`);
	buildSearchParams(url.searchParams, params);

	const res = await window.fetch(url, {
		method: 'GET',
		credentials: 'omit',
		headers: new Headers({
			"Content-Type": "application/json",
		}),
	});

	return await res.json();
}

export function useApi() {
	return {
		post,
		get,
	}
}
