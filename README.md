# forum-authentication

This project is a web forum that allows users to communicate, share posts, comment, and interact with one another through likes/dislikes, filtering, and more.

---

## Objectives

- **User Communication**:  
  Allow users to create posts and comments to facilitate discussion.

- **Categorized Posts**:  
  Users can associate one or more categories to their posts, functioning similarly to subforums dedicated to specific topics.

- **Likes and Dislikes**:  
  Registered users can like or dislike posts and comments. The total counts of likes and dislikes will be visible to all users.

- **Filtering**:  
  Implement filtering for posts by:

  - Categories
  - Created posts (for the logged-in user)
  - Liked posts (for the logged-in user)

- **Authentication**:  
   A user can create an account by using the following third-party services:
  - Google
  - GitHub

---

## Technologies Used

- **Language**: Go (Golang), HTML, CSS

- **Database**: SQLite

  - SQLite is chosen for its simplicity as an embedded database and ease of integration in web applications.

- **Authentication and Session Management**:

  - User registration and login with email, username, and password.
  - Use cookies for session management with a 24-hour time period.
    - Encrypting passwords using `bcrypt`.
    - Implementing session identifiers using `UUID`.

- **Docker**:
  - Containerizing the application for consistent deployment and easy environment management.

---

## Authentication

### User Registration

- **Input Requirements:**

  - **Email**: Must be unique. Cannot register a user if the email is already registered.

  - **Username**

  - **Password**: Encrypted when stored (uses `bcrypt` for encryption).

### Login

- Validate user credentials against stored records.
- Check that the password provided matches the encrypted password in the database.
- On successful login, it creates a session cookie with an expiration date; with only one active session per user.

---

### Communication

- **Posts & Comments:**
  - Only registered users can create posts and comments.
  - Posts can be associated with one or more categories.
  - Both posts and comments are visible to all users, regardless of registration status.
  - Non-registered users can only view posts and comments but cannot interact with them (no reaction; like, dislike, or comments).

---

### Likes and Dislikes

- **Functionality:**
  - Only registered users can like or dislike posts and comments.
  - The count of likes and dislikes is visible to all users.

---

### Filtering

- **Categories:**  
  Users can filter posts by specific categories (similar to subforums).

- **Created Posts:**  
  Registered users can filter posts that they have created.
- **Liked Posts:**  
  Registered users can filter posts that they have liked.

---

## Installation

1. Clone the repository:

   ```bash
   git clone https://learn.zone01kisumu.ke/git/johnodhiambo0/forum-authentication.git
   cd forum-authentication
   ```

2. Create a `.env` file in the root directory and add the following credentials:

   ```env
   GOOGLE_CLIENT_ID=
   GOOGLE_CLIENT_SECRET=

   GITHUB_CLIENT_ID=
   GITHUB_CLIENT_SECRET=
   ```

   Replace the values with your own credentials obtained from Google and GitHub.

3. Compile and run the program with a file as input:

   ```bash
   go run main.go
   ```

### Setting up Google and GitHub OAuth

#### Google OAuth Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/).
2. Create a new project or select an existing one.
3. Navigate to **APIs & Services > Credentials**.
4. Click **Create Credentials** > **OAuth 2.0 Client ID**.
5. Configure the OAuth consent screen with the following information:
   - **Authorized JavaScript origins:** `http://localhost:9000`
   - **Authorized redirect URIs:** `http://localhost:9000/auth/google/callback`
   - **Authorized redirect URIs:** `http://localhost:9000/auth/google/signin/callback`
6. Copy the **Client ID** and **Client Secret** and paste them into the `.env` file.

#### GitHub OAuth Setup

1. Go to [GitHub Developer Settings](https://github.com/settings/developers).
2. Click **New OAuth App**.
3. Fill in the application details:
   - **Homepage URL**: `http://localhost:9000`
   - **Authorization callback URL**: `http://localhost:9000/auth/github/callback`
   - **Authorization callback URL**: `http://localhost:9000/auth/github/signin/callback`
4. Register the application.
5. Copy the **Client ID** and **Client Secret** and paste them into the `.env` file.

### Docker

To ensure ease of deployment and consistency across environments, this project uses Docker.

**Building an Image**:

```bash
docker build -t forum .
```

- You can build using `docker-compose.yml`:

```bash
docker compose up --build
```

## Contribution

- To make a contribution to the project, open an issue with a title, a tag, and a description of your idea on the [repository issues' page](https://github.com/jesee-kuya/forum/issues).

## License

This project is licensed under [MIT](https://github.com/jesee-kuya/forum/blob/main/LICENSE).
