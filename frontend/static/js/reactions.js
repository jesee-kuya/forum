document.querySelectorAll('.like-button, .dislike-button, .like-comment-button, .dislike-comment-button').forEach((button) => {
  button.addEventListener('click', function (event) {
    event.preventDefault();

    const postId = button.getAttribute('data-posted-id');
    const reaction = button.getAttribute('data-reaction');

    fetch('/reaction', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      credentials: 'include',
      body: `post_id=${postId}&reaction=${reaction}`,
    })
      .then((response) => response.text())
      .then((data) => {
        window.location.reload();
      })
      .catch((err) => console.error(err));
  });
});