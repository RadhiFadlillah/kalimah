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
