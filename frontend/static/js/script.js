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

  createPostBtn.addEventListener('click', () => {
    createPostSection.classList.toggle('hidden');
  });
});

// Uses fetch API to make asynchronous requests to the backend when the like button is clicked.
document.querySelectorAll('.like-button').forEach((button) => {
  button.addEventListener('click', async function () {
    console.log('Like button clicked');

    const postId = this.getAttribute('data-post-id');
    if (!postId) return;

    const response = await fetch('/like', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ post_id: parseInt(postId) }),
    });

    if (response.ok) {
      const data = await response.json();
      // Update the like count dynamically
      this.querySelector('span').textContent = data.newLikeCount;
    } else {
      console.error('Failed to update like count');
    }
  });
});

// Uses fetch API to make asynchronous requests to the backend when the comment button is clicked.
document.querySelectorAll('.submit-comment').forEach((button) => {
  button.addEventListener('click', async function () {
    const postId = this.dataset.postId;
    const commentText = this.previousElementSibling.value;

    const response = await fetch('/comment', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ post_id: postId, comment: commentText }),
    });

    if (response.ok) {
      const data = await response.json();
      const commentSection = this.closest('.comments-section');
      commentSection.innerHTML += `<p><strong>You:</strong> ${data.comment}</p>`;
    }
  });
});
