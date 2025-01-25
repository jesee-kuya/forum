'use strict';

const themeToggler = document.querySelector('.theme-toggler');

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

themeToggler.addEventListener('click', toggleTheme);

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

document.addEventListener('DOMContentLoaded', () => {
  const createPostSection = document.querySelector('.create-post');
  const createPostBtn = document.querySelector('.floating-create-post-btn');

  createPostBtn.addEventListener('click', () => {
    createPostSection.classList.toggle('hidden');
  });
});
