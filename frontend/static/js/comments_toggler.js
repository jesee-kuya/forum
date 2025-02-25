document.addEventListener('DOMContentLoaded', () => {
  document.querySelectorAll('.comment-button').forEach((button) => {
    const post = button.closest('.post');
    const commentsSection = post?.querySelector('.comments-section');
    const span = button.querySelector('.comment-count');

    if (!commentsSection || !span) return;

    const postId = post.getAttribute('data-post-id');
    const commentCount = parseInt(span.textContent.trim(), 10);
    const wasOpened =
      localStorage.getItem(`comments-visible-${postId}`) === 'true';

    // Always hide comments initially
    commentsSection.classList.add('hidden');
    commentsSection.style.display = 'none';

    // If count >= 1 and was previously opened, show it
    if (commentCount > 0 && wasOpened) {
      commentsSection.classList.remove('hidden');
      commentsSection.style.display = 'block';
    }

    // Toggle comments when clicking the button
    button.addEventListener('click', () => {
      const isNowVisible = commentsSection.classList.toggle('hidden');
      commentsSection.style.display = isNowVisible ? 'none' : 'block';

      // Save state in localStorage
      localStorage.setItem(`comments-visible-${postId}`, !isNowVisible);
    });
  });
});
