document.addEventListener('DOMContentLoaded', function () {
  const signinForm = document.getElementById('signin-form');
  const popup = document.getElementById('message-popup');

  function showMessage(message, isSuccess) {
    popup.textContent = message;
    popup.classList.remove('success', 'error');

    popup.classList.add('show', isSuccess ? 'success' : 'error');

    setTimeout(() => {
      popup.classList.remove('show', 'success', 'error');
    }, 3000);
  }

  const urlParams = new URLSearchParams(window.location.search);
  if (urlParams.has('status') && urlParams.get('status') === 'success') {
    showMessage('Sign Up Successful!', true);
    history.replaceState(null, '', window.location.pathname);
  } else if (urlParams.has('error')) {
    showMessage(`OAuth Signup failed: ${urlParams.get('error')}`, false);
    history.replaceState(null, '', window.location.pathname);
  }

  // Attach event listeners to OAuth buttons
  document.querySelector('.google-btn').addEventListener('click', (e) => {
    e.preventDefault();
    window.location.href = '/auth/google/signin';
  });

  document.querySelector('.github-btn').addEventListener('click', (e) => {
    e.preventDefault();
    window.location.href = '/auth/github/signin';
  });

  signinForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const signinFormData = new URLSearchParams(new FormData(signinForm));
    console.log(Object.fromEntries(signinFormData));

    try {
      const response = await fetch('/sign-in', {
        method: 'POST',
        body: signinFormData,
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      });

      const data = await response.json();

      if (data.success) {
        showMessage('Sign In Successful!', true);

        setTimeout(() => {
          window.location.href = '/';
        }, 1000);
      } else {
        showMessage('Operation failed. Please check your input.', false);
      }
    } catch (error) {
      console.error('Error:', error);
      showMessage('Operation failed. Please check your input.', false);
    }
  });
});
