package auth

// OAuthSuccessHTML is the HTML page shown after successful authentication
const OAuthSuccessHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Authentication Successful - Ubuntu Software</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
      background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
    }
    .container {
      text-align: center;
      padding: 3rem;
      background: rgba(255, 255, 255, 0.05);
      border-radius: 16px;
      backdrop-filter: blur(10px);
      border: 1px solid rgba(255, 255, 255, 0.1);
      max-width: 480px;
    }
    .logo { width: 180px; height: auto; margin-bottom: 2rem; }
    .check {
      width: 80px; height: 80px;
      background: #22c55e;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      margin: 0 auto 1.5rem;
    }
    .check svg { width: 40px; height: 40px; }
    h1 { font-size: 1.75rem; font-weight: 600; margin-bottom: 0.75rem; }
    p { color: rgba(255, 255, 255, 0.7); font-size: 1rem; line-height: 1.5; }
    .hint { margin-top: 2rem; font-size: 0.875rem; color: rgba(255, 255, 255, 0.5); }
  </style>
</head>
<body>
  <div class="container">
    <img src="https://www.ubuntusoftware.net/images/logo.svg" alt="Ubuntu Software" class="logo">
    <div class="check">
      <svg fill="none" stroke="white" stroke-width="3" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
      </svg>
    </div>
    <h1>Authentication Successful</h1>
    <p>Your Google Cloud credentials have been configured. You can now use Terraform and other Google Cloud tools.</p>
    <p class="hint">You may close this window.</p>
  </div>
</body>
</html>`

// OAuthErrorHTML is the HTML page shown after failed authentication
// Use with fmt.Sprintf to inject the error message
const OAuthErrorHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Authentication Failed - Ubuntu Software</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
      background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
    }
    .container {
      text-align: center;
      padding: 3rem;
      background: rgba(255, 255, 255, 0.05);
      border-radius: 16px;
      backdrop-filter: blur(10px);
      border: 1px solid rgba(255, 255, 255, 0.1);
      max-width: 480px;
    }
    .logo { width: 180px; height: auto; margin-bottom: 2rem; }
    .error-icon {
      width: 80px; height: 80px;
      background: #ef4444;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      margin: 0 auto 1.5rem;
    }
    .error-icon svg { width: 40px; height: 40px; }
    h1 { font-size: 1.75rem; font-weight: 600; margin-bottom: 0.75rem; }
    p { color: rgba(255, 255, 255, 0.7); font-size: 1rem; line-height: 1.5; }
    .error-msg {
      margin-top: 1rem;
      padding: 1rem;
      background: rgba(239, 68, 68, 0.2);
      border-radius: 8px;
      font-family: monospace;
      font-size: 0.875rem;
    }
    .hint { margin-top: 2rem; font-size: 0.875rem; color: rgba(255, 255, 255, 0.5); }
  </style>
</head>
<body>
  <div class="container">
    <img src="https://www.ubuntusoftware.net/images/logo.svg" alt="Ubuntu Software" class="logo">
    <div class="error-icon">
      <svg fill="none" stroke="white" stroke-width="3" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
      </svg>
    </div>
    <h1>Authentication Failed</h1>
    <p>There was a problem authenticating with Google Cloud.</p>
    <div class="error-msg">%s</div>
    <p class="hint">Please try again or check your Google Cloud configuration.</p>
  </div>
</body>
</html>`
