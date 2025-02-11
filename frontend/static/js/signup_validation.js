document.addEventListener('DOMContentLoaded', () => {
  const nameInput = document.getElementById('username');
  const emailInput = document.getElementById('email');
  const passwordInput = document.getElementById('password');
  const confirmPasswordInput = document.getElementById('confirmed-password');
  const signupForm = document.getElementById('signup-form');

  // Feedback elements
  const nameFeedback = document.createElement('p');
  nameFeedback.className = 'feedback-message';
  nameInput.parentNode.appendChild(nameFeedback);

  const emailFeedback = document.createElement('p');
  emailFeedback.className = 'feedback-message';
  emailInput.parentNode.appendChild(emailFeedback);

  // Prevent excessive calls using debounce
  function debounce(func, delay) {
    let timeout;
    return (...args) => {
      clearTimeout(timeout);
      timeout = setTimeout(() => func(...args), delay);
    };
  }

  // Check credentials availability
  async function checkAvailability(field, value, feedbackElement) {
    if (!value.trim()) {
      feedbackElement.textContent = '';
      return;
    }

    try {
      const response = await fetch(
        `/validate?${field}=${encodeURIComponent(value)}`
      );
      const data = await response.json();

      if (data.available) {
        feedbackElement.textContent = `${
          field.charAt(0).toUpperCase() + field.slice(1)
        } is available`;
        feedbackElement.style.color = 'green';
      } else {
        feedbackElement.textContent = `${
          field.charAt(0).toUpperCase() + field.slice(1)
        } is taken`;
        feedbackElement.style.color = 'red';
      }
    } catch (error) {
      console.error('Error validating input:', error);
    }
  }

  nameInput.addEventListener(
    'input',
    debounce(
      () => checkAvailability('username', nameInput.value, nameFeedback),
      1000
    )
  );
  emailInput.addEventListener(
    'input',
    debounce(
      () => checkAvailability('email', emailInput.value, emailFeedback),
      1000
    )
  );

  function validatePasswordStength(password) {
    if (password.length < 8) return 'Must be at least 8 characters.';
    if (!/[A-Z]/.test(password))
      return 'Include at least one uppercase letter.';
    if (!/[a-z]/.test(password))
      return 'Include at least one lowercase letter.';
    if (!/[0-9]/.test(password)) return 'Include at least one number.';
    if (!/[!@#$%^&*]/.test(password))
      return 'Include at least one special character.';
    return '';
  }

  // Show password strength validation
  passwordInput.addEventListener('input', () => {
    const passwordError = validatePasswordStength(passwordInput.value);
    passwordInput.setCustomValidity(passwordError);
    passwordInput.reportValidity();
  });

  // Confirm validation
  confirmPasswordInput.addEventListener('input', () => {
    if (passwordInput.value !== confirmPasswordInput.value) {
      confirmPasswordInput.setCustomValidity('Passwords do not match.');
    } else {
      confirmPasswordInput.setCustomValidity('');
    }
    confirmPasswordInput.reportValidity();
  });

  // Prevent submission of validation fails
  signupForm.addEventListener('submit', (e) => {
    if (!signupForm.checkValidity()) {
      e.preventDefault();
    }
  });

  const popup = document.getElementById('message-popup');
  signupForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const signUpFormData = new URLSearchParams(new FormData(signupForm));
    console.log(Object.fromEntries(signUpFormData));

    try {
      const response = await fetch('/sign-up', {
        method: 'POST',
        body: signUpFormData,
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      });

      const data = await response.json();

      if (data.success) {
        showMessage('Sign Up Successful!', true);
        setTimeout(() => {
          window.location.href = '/sign-in';
        }, 1000);
      } else {
        showMessage('Sign Up Failed. Please check your input.', false);
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
