document.addEventListener('DOMContentLoaded', function () {
  const form = document.getElementById('signup-form');
  const errorMessage = document.querySelector('.error-message');

  form.addEventListener('submit', async function (event) {
    event.preventDefault();

    // Convert FormData to URL-encoded format
    const formData = new URLSearchParams(new FormData(form));
    const response = await fetch('/sign-up', {
      method: 'POST',
      body: formData,
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    });

    if (response.redirected) {
      window.location.href = response.url;
      return;
    }

    const result = await response.json();

    if (!result.success) {
      // Show error message if validation fails
      errorMessage.classList.add('show');
      setTimeout(() => {
        errorMessage.classList.remove('show');
      }, 3000);
    }
  });
});
