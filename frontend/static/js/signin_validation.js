document.addEventListener('DOMContentLoaded', function () {
  const signInForm = document.getElementById('signin-form');
  const errorMessage = document.querySelector('.error-message');

  signInForm.addEventListener('submit', async function (e) {
    e.preventDefault();

    // Convert FormData to URL-encoded format
    const signInFormData = new URLSearchParams(new FormData(signInForm));
    const signInResponse = await fetch('/sign-in', {
      method: 'POST',
      body: signInFormData,
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    });

    if (signInResponse.redirected) {
      window.location.href = signInResponse.url;
      return;
    }

    const signInResult = await signInResponse.json();
    if (!signInResult.success) {
      errorMessage.classList.add('show');
      setTimeout(() => {
        errorMessage.classList.remove('show');
      }, 3000);
    }
  });
});
