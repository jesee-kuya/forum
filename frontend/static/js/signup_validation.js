document.addEventListener('DOMContentLoaded', function () {
  const signUpForm = document.getElementById('signup-form');
  const errorMessage = document.querySelector('.error-message');

  signUpForm.addEventListener('submit', async function (event) {
    event.preventDefault();

    // Convert FormData to URL-encoded format
    const signUpFormData = new URLSearchParams(new FormData(signUpForm));
    const signUpResponse = await fetch('/sign-up', {
      method: 'POST',
      body: signUpFormData,
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    });

    if (signUpResponse.redirected) {
      window.location.href = signUpResponse.url;
      return;
    }

    const signUpResult = await signUpResponse.json();
    if (!signUpResult.success) {
      // Show error message if validation fails
      errorMessage.classList.add('show');
      setTimeout(() => {
        errorMessage.classList.remove('show');
      }, 3000);
    }
  });
});
