const handleOAuthResponse = () => {
  const urlParams = new URLSearchParams(window.location.search);
  const error = urlParams.get('error');
  const success = urlParams.get('success');

  if (error) {
    const errorMessages = {
      invalid_state: 'Authentication failed: Invalid state',
      token_exchange_failed:
        'Authentication failed: Unable to connect to service',
      user_info_failed: 'Unable to get user information',
      auth_failed: 'Authentication failed',
      session_error: 'Session creation failed',
      no_account: 'No account found. Please sign up first.',
      default: 'Authentication failed. Please try again.',
    };

    showMessage(errorMessages[error] || errorMessages.default, false);
  } else if (success) {
    showMessage('Authentication successful!', true);
    // Redirect after showing success message
    setTimeout(() => {
      window.location.href = success === 'signup' ? '/sign-in' : '/home';
    }, 1500);
  }
};

function showMessage(message, isSuccess) {
  const popup = document.getElementById('message-popup');
  popup.textContent = message;
  popup.className = `message-popup show ${isSuccess ? 'success' : 'error'}`;

  setTimeout(() => {
    popup.classList.remove('show');
  }, 3000);
}
