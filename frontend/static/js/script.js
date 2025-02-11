// Apply theme based on the value in localStorage
const applyTheme = (theme) => {
  if (theme === 'dark') {
    document.body.classList.add('dark-theme');
  } else {
    document.body.classList.remove('dark-theme');
  }
};

const toggleTheme = () => {
  const currentTheme = document.body.classList.contains('dark-theme')
    ? 'dark'
    : 'light';
  const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
  applyTheme(newTheme);
  localStorage.setItem('theme', newTheme);
};
const savedTheme = localStorage.getItem('theme') || 'light';
applyTheme(savedTheme);

document.querySelector('.theme-toggler').addEventListener('click', toggleTheme);

// Toggle comment button
document.querySelectorAll('.comment-button').forEach((button) => {
  button.addEventListener('click', () => {
    const post = button.closest('.post');
    const commentsSection = post.querySelector('.comments-section');
    commentsSection.classList.toggle('hidden');
    commentsSection.style.display = commentsSection.classList.contains('hidden')
      ? 'none'
      : 'block';
  });
});

// Toggle post creation window
document.addEventListener('DOMContentLoaded', () => {
  const createPostSection = document.querySelector('.create-post');
  const createPostBtn = document.querySelector('.floating-create-post-btn');
  const postsContainer = document.querySelector('main.posts');

  createPostBtn.addEventListener('click', () => {
    createPostSection.classList.toggle('hidden');

    if (postsContainer) {
      postsContainer.scrollTo({
        top: 0,
        behavior: 'smooth',
      });
    }
  });
});

// Password toggle
document.querySelectorAll('.toggle-password').forEach((button) => {
  button.addEventListener('click', () => {
    const input = document.getElementById(button.dataset.target);
    if (input.type === 'password') {
      input.type = 'text';
    } else {
      input.type = 'password';
    }
  });
});
