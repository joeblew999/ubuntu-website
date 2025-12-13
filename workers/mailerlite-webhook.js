/**
 * Cloudflare Worker: Web3Forms → MailerLite Webhook Handler
 *
 * This worker receives webhook POSTs from Web3Forms and adds subscribers
 * to MailerLite. Deploy to Cloudflare Workers for serverless operation.
 *
 * Environment Variables (set via wrangler.toml or dashboard):
 *   MAILERLITE_API_KEY   - Your MailerLite API key (secret)
 *   MAILERLITE_GROUP_ID  - Group ID to assign subscribers to
 *
 * Endpoints:
 *   POST /webhook  - Receive Web3Forms submission
 *   GET  /health   - Health check
 *   GET  /         - Info page
 */

const MAILERLITE_API = 'https://connect.mailerlite.com/api';

export default {
	async fetch(request, env, ctx) {
		const url = new URL(request.url);

		// Route handling
		switch (url.pathname) {
			case '/':
				return handleInfo();
			case '/health':
				return handleHealth();
			case '/webhook':
				return handleWebhook(request, env);
			default:
				return new Response('Not Found', { status: 404 });
		}
	},
};

/**
 * Info page - shows worker status
 */
function handleInfo() {
	return new Response(
		JSON.stringify({
			name: 'mailerlite-webhook',
			version: '1.0.0',
			endpoints: {
				'/': 'Info (this page)',
				'/health': 'Health check',
				'/webhook': 'POST - Receive Web3Forms submission',
			},
		}),
		{
			headers: { 'Content-Type': 'application/json' },
		}
	);
}

/**
 * Health check endpoint
 */
function handleHealth() {
	return new Response('OK', { status: 200 });
}

/**
 * Main webhook handler - receives Web3Forms submissions
 */
async function handleWebhook(request, env) {
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
