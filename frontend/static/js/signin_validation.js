document.addEventListener('DOMContentLoaded', function () {
  const signinForm = document.getElementById('signin-form');
  const popup = document.getElementById('message-popup');
  
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
        showMessage('Sign In Failed. Please check your input.', false);
      }
    } catch (error) {
      console.error('Error:', error);
      showMessage('An error occurred. Try again later.', false);
    }
  });

  function showMessage(message, isSuccess) {
    popup.textContent = message;
    popup.classList.add('show');

    setTimeout(() => {
      popup.classList.remove('show');
    }, 3000);
  }
});
