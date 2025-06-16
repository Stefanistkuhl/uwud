document.addEventListener('DOMContentLoaded', () => {
	const cards = document.querySelectorAll('.card');

	cards.forEach(card => {
		const healthCheckUrl = card.dataset.healthCheckUrl;
		const healthIcon = card.querySelector('.health-icon');

		if (healthCheckUrl) {
			checkHealth(healthCheckUrl, healthIcon);
		} else {
			healthIcon.innerHTML = '';
			healthIcon.title = '';
			healthIcon.classList.remove('pending', 'success', 'failure');
		}
	});

	function svgPending() {
		return `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor"
			stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-label="Checking...">
			<circle cx="12" cy="12" r="10" stroke-opacity="0.3"/>
			<path d="M12 6v6l4 2" />
			</svg>`;
	}
	function svgUp() {
		return `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#22c55e"
			stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" aria-label="Healthy">
			<circle cx="12" cy="12" r="10" stroke="#22c55e" stroke-opacity="0.2" fill="none"/>
			<polyline points="8 12.5 11 15.5 16 9.5" fill="none"/>
			</svg>`;
	}
	function svgDown() {
		return `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#ef4444"
			stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" aria-label="Unhealthy">
			<circle cx="12" cy="12" r="10" stroke="#ef4444" stroke-opacity="0.2" fill="none"/>
			<line x1="8" y1="8" x2="16" y2="16"/>
			<line x1="16" y1="8" x2="8" y2="16"/>
			</svg>`;
	}

	function checkHealth(url, iconElement) {
		iconElement.innerHTML = svgPending();
		iconElement.title = 'Checking...';
		iconElement.classList.remove('success', 'failure');
		iconElement.classList.add('pending');
		fetch(`/healthcheck?url=${encodeURIComponent(url)}`, { method: 'GET' })
			.then(response =>
				response.json().then(data => ({ status: response.status, data }))
			)
			.then(({ status, data }) => {
				if (status === 200) {
					iconElement.innerHTML = svgUp();
					iconElement.title = 'Healthy';
					iconElement.classList.remove('pending', 'failure');
					iconElement.classList.add('success');
				} else {
					iconElement.innerHTML = svgDown();
					iconElement.title = `Unhealthy (Status code: ${data.code || status})`;
					iconElement.classList.remove('pending', 'success');
					iconElement.classList.add('failure');
				}
			})
			.catch(error => {
				iconElement.innerHTML = svgDown();
				iconElement.title = 'Health check failed';
				iconElement.classList.remove('pending', 'success');
				iconElement.classList.add('failure');
				console.error(`Health check failed for ${url}:`, error);
			});
	}
});
