export async function getRequest(url: string): Promise<any> {
	// Send GET request
	const resp = await fetch(url);

	// Check for error message
	if (!resp.ok) {
		let responseText = await resp.text();
		throw Error(`${responseText.trim()} (${resp.status})`);
	}

	// Return body
	if (resp.headers.get('content-type') === 'application/json') {
		return await resp.json();
	}

	return await resp.text();
}

export async function postRequest(url: string, body?: any): Promise<any> {
	// Prepare request
	let request: Record<string, any> = {
		method: 'post',
		headers: {},
	};

	if (body != null) {
		request.body = JSON.stringify(body);
		request.headers['Content-Type'] = 'application/json; charset=utf-8';
	}

	// Send POST request
	const resp = await fetch(url, request);

	// Check for error message
	if (!resp.ok) {
		let responseText = await resp.text();
		throw Error(`${responseText.trim()} (${resp.status})`);
	}

	// Return body
	if (resp.headers.get('content-type') === 'application/json') {
		return await resp.json();
	}

	return await resp.text();
}
