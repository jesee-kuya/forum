function openOAuthPopup(url, name = 'OAuth Popup', width = 600, height = 600) {
  const left = (window.innerWidth - width) / 2 + window.screenX;
  const top = (window.innerHeight - height) / 2 + window.screenY;

  const popup = window.open(
    url,
    name,
    `width=${width},height=${height},left=${left},top=${top},resizable=yes,scrollbars=yes,status=yes`
  );

  // Create a promise that resolves when the popup completes
  return new Promise((resolve, reject) => {
    // Check if popup was blocked
    if (!popup || popup.closed || typeof popup.closed === 'undefined') {
      reject(new Error('Popup blocked. Please allow popups for this site.'));
      return;
    }

    // Poll to check when the popup is closed
    const pollTimer = setInterval(() => {
      try {
        // If popup redirected to same origin, we can access location
        if (popup.location.href.indexOf(window.location.origin) !== -1) {
          clearInterval(pollTimer);

          // Extract URL parameters
          const params = new URLSearchParams(popup.location.search);
          popup.close();

          if (params.has('status') && params.get('status') === 'success') {
            resolve({ success: true });
          } else if (params.has('error')) {
            resolve({ success: false, error: params.get('error') });
          } else {
            resolve({ success: true });
          }
        }
      } catch (e) {
        // Cross-origin error, ignore as expected during OAuth process
      }

      // Check if popup was closed
      if (popup.closed) {
        clearInterval(pollTimer);
        resolve({ success: false, error: 'Authentication canceled' });
      }
    }, 500);

    // Add window message event listener for communication
    window.addEventListener('message', function receiveMessage(event) {
      // Verify origin to prevent XSS attacks
      if (event.origin !== window.location.origin) return;

      if (event.data.type === 'oauth_complete') {
        clearInterval(pollTimer);
        window.removeEventListener('message', receiveMessage);
        popup.close();
        resolve(event.data);
      }
    });
  });
}
