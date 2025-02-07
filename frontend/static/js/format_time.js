document.querySelectorAll('.post-time time').forEach((el) => {
  const rawTimestamp = el.getAttribute('datetime');
  const parsedTimestamp = new Date(rawTimestamp.replace(' +0000 UTC', 'Z'));

  if (!isNaN(parsedTimestamp)) {
    el.innerText = formatTimestamp(parsedTimestamp);
  } else {
    console.warn(`Invalid date format: ${rawTimestamp}`);
  }
});

function formatTimestamp(timestamp) {
  const now = new Date();
  const pastTime = new Date(timestamp);
  const timeDifference = Math.floor((now - pastTime) / 1000);

  if (timeDifference < 60) {
    return timeDifference === 1
      ? '1 second ago'
      : `${timeDifference} seconds ago`;
  } else if (timeDifference < 3600) {
    return Math.floor(timeDifference / 60) === 1
      ? '1 minute ago'
      : `${Math.floor(timeDifference / 60)} minutes ago`;
  } else if (timeDifference < 86400) {
    return Math.floor(timeDifference / 3600) === 1
      ? '1 hour ago'
      : `${Math.floor(timeDifference / 3600)} hours ago`;
  } else if (timeDifference < 604800) {
    return Math.floor(timeDifference / 86400) === 1
      ? '1 day ago'
      : `${Math.floor(timeDifference / 86400)} days ago`;
  } else if (timeDifference < 2592000) {
    return Math.floor(timeDifference / 604800) === 1
      ? '1 week ago'
      : `${Math.floor(timeDifference / 604800)} weeks ago`;
  } else if (timeDifference < 31536000) {
    return Math.floor(timeDifference / 2592000) === 1
      ? '1 month ago'
      : `${Math.floor(timeDifference / 2592000)} months ago`;
  } else {
    return Math.floor(timeDifference / 31536000) === 1
      ? '1 year ago'
      : `${Math.floor(timeDifference / 31536000)} years ago`;
  }
}

fetch('/posts')
  .then((response) => response.json())
  .then((posts) => {
    posts.forEach((post) => {
      const timeElement = document.querySelector(
        `.post-time time[datetime="${post.created_on}"]`
      );

      if (timeElement) {
        const timestamp = new Date(post.created_on);
        timeElement.innerText = formatTimestamp(timestamp);
      } else {
        console.warn(
          `No matching <time> element found for timestamp: ${post.created_on}`
        );
      }
    });
  })
  .catch((error) => console.error('Error fetching posts:', error));
