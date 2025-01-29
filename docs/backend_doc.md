## Authentication Endpoints

### 1. **Register User**
**POST** `/api/register`

**Description:** Create a new user account.

**Request Body:**
```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

**Response:**
- **201 Created:**
```json
{
  "message": "User registered successfully"
}
```
- **400 Bad Request:** Invalid input.
- **500 Internal Server Error:** User could not be registered.


### 2. **Login**
**POST** `/api/login`

**Description:** Authenticate a user and return a session token.

**Request Body:**
```json
{
  "email": "string",
  "password": "string"
}
```

**Response:**
- **200 OK:**
```json
{
  "message": "Login successful",
  "token": "string"
}
```
- **401 Unauthorized:** Invalid credentials.
- **500 Internal Server Error:** Login failed.


### 3. **Logout**
**POST** `/api/logout`

**Description:** Log out the user and invalidate the session token.

**Headers:**
- `Authorization: Bearer <token>`

**Response:**
- **200 OK:**
```json
{
  "message": "Logout successful"
}
```
- **401 Unauthorized:** Invalid or missing token.


## Posts Endpoints

### 4. **Get All Posts**
**GET** `/api/posts`

**Description:** Retrieve all posts in the forum.

**Response:**
- **200 OK:**
```json
[
  {
    "id": 1,
    "user_id": 2,
    "title": "string",
    "body": "string",
    "created_on": "timestamp",
    "post_type": "string"
  }
]
```
- **500 Internal Server Error:** Could not fetch posts.


### 5. **Create a Post**
**POST** `/api/posts`

**Description:** Create a new post.

**Headers:**
- `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "title": "string",
  "body": "string",
  "post_type": "string"
}
```

**Response:**
- **201 Created:**
```json
{
  "message": "Post created successfully"
}
```
- **400 Bad Request:** Invalid input.
- **401 Unauthorized:** Missing or invalid token.
- **500 Internal Server Error:** Could not create the post.


### 6. **Like a Post**
**POST** `/api/posts/like`

**Description:** Like a specific post.

**Headers:**
- `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "post_id": 1,
  "reaction": "like"
}
```

**Response:**
- **200 OK:**
```json
{
  "message": "Post liked successfully"
}
```
- **400 Bad Request:** Invalid input.
- **401 Unauthorized:** Missing or invalid token.
- **500 Internal Server Error:** Could not like the post.


### 7. **Comment on a Post**
**POST** `/api/posts/comment`

**Description:** Add a comment to a post.

**Headers:**
- `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "post_id": 1,
  "body": "string"
}
```

**Response:**
- **201 Created:**
```json
{
  "message": "Comment added successfully"
}
```
- **400 Bad Request:** Invalid input.
- **401 Unauthorized:** Missing or invalid token.
- **500 Internal Server Error:** Could not add the comment.


## Categories Endpoints

### 8. **Get All Categories**
**GET** `/api/categories`

**Description:** Retrieve all categories available.

**Response:**
- **200 OK:**
```json
[
  {
    "id": 1,
    "name": "string"
  }
]
```
- **500 Internal Server Error:** Could not fetch categories.


### 9. **Filter Posts by Category**
**GET** `/api/posts/filter?category=<name>`

**Description:** Retrieve posts filtered by a specific category.

**Response:**
- **200 OK:**
```json
[
  {
    "id": 1,
    "user_id": 2,
    "title": "string",
    "body": "string",
    "created_on": "timestamp",
    "post_type": "string"
  }
]
```
- **500 Internal Server Error:** Could not fetch filtered posts.


## Users Endpoints

### 10. **Get User Profile**
**GET** `/api/users/<id>`

**Description:** Retrieve details of a specific user by ID.

**Response:**
- **200 OK:**
```json
{
  "id": 1,
  "username": "string",
  "email": "string",
  "joined_on": "timestamp"
}
```
- **404 Not Found:** User does not exist.
- **500 Internal Server Error:** Could not fetch user details.


### 11. **Delete a User**
**DELETE** `/api/users/<id>`

**Description:** Delete a specific user by ID.

**Headers:**
- `Authorization: Bearer <token>`

**Response:**
- **200 OK:**
```json
{
  "message": "User deleted successfully"
}
```
- **401 Unauthorized:** Missing or invalid token.
- **403 Forbidden:** User is not authorized to perform this action.
- **500 Internal Server Error:** Could not delete the user.


## Error Handling

All endpoints return consistent error messages in the following format:
```json
{
  "error": "string",
  "status": 400
}
```
