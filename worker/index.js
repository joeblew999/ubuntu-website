/**
 * The site's Worker: the Hugo build's pages, from its static assets (public/, wrangler.jsonc's assets), and under
 * /api/ the routes that need code. Only /api/* runs this script (assets.run_worker_first); every other request is an
 * asset, served by Cloudflare without it.
 *
 *   GET  /api/health              health check
 *   POST /api/webhook/mailerlite  Web3Forms' submission → a MailerLite subscriber
 */

const MAILERLITE_API = 'https://connect.mailerlite.com/api';

export default {
	async fetch(request, env) {
		const url = new URL(request.url);
		switch (url.pathname) {
			case '/api/health':
				return new Response('OK', { status: 200 });
			case '/api/webhook/mailerlite':
				return handleMailerLite(request, env);
			default:
				return env.ASSETS.fetch(request);
		}
	},
};

/**
 * Web3Forms' submission, POSTed to /api/webhook/mailerlite: its email (and name, company) added to MailerLite, in
 * MAILERLITE_GROUP_ID's group when there is one. MAILERLITE_API_KEY is the Worker's secret.
 */
async function handleMailerLite(request, env) {
	// Only accept POST
	if (request.method !== 'POST') {
		return new Response('Method not allowed', { status: 405 });
	}

	// Check for API key
	if (!env.MAILERLITE_API_KEY) {
		console.error('MAILERLITE_API_KEY not configured');
		return new Response('Server configuration error', { status: 500 });
	}

	try {
		// Parse form data (Web3Forms sends form-urlencoded)
		const contentType = request.headers.get('Content-Type') || '';
		let data;

		if (contentType.includes('application/json')) {
			data = await request.json();
		} else {
			// Form-encoded (default from Web3Forms)
			const formData = await request.formData();
			data = Object.fromEntries(formData.entries());
		}

		// Extract fields
		const email = data.email;
		const name = data.name || '';
		const company = data.company || '';
		const platform = data.platform || '';
		const industry = data.industry || '';
		const usecase = data.usecase || '';

		// Validate email
		if (!email) {
			console.log('Webhook received without email');
			return new Response('Email required', { status: 400 });
		}

		console.log(`New submission: ${email}`);

		// Build subscriber data
		const subscriberData = {
			email,
			fields: {
				name,
				company,
			},
		};

		// Add to group if configured
		if (env.MAILERLITE_GROUP_ID) {
			subscriberData.groups = [env.MAILERLITE_GROUP_ID];
		}

		// Call MailerLite API
		const response = await fetch(`${MAILERLITE_API}/subscribers`, {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${env.MAILERLITE_API_KEY}`,
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(subscriberData),
		});

		if (!response.ok) {
			const error = await response.text();
			console.error(`MailerLite API error: ${response.status} - ${error}`);
			// Still return OK to Web3Forms (prevent retries)
			return new Response('Received (MailerLite error logged)', { status: 200 });
		}

		const result = await response.json();
		console.log(`Subscriber added: ${result.data.id}`);

		return new Response('OK', { status: 200 });
	} catch (error) {
		console.error('Webhook error:', error);
		// Return OK to prevent Web3Forms retries
		return new Response('Received (error logged)', { status: 200 });
	}
}
