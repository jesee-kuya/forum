# API Documentation for Forum Backend

## Base URL

```
http://localhost:8080
```

## Authentication

### Register User

**Endpoint:**

```
POST /api/register
```

**Request Body:**

```json
{
  "username": "kherld",
  "email": "kherld@forum.com",
  "password": "frmpwd"
}
```

**Test with cURL:**

```
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username": "kherld", "email": "kherld@forum.com", "password": "frmpwd"}'
```

### Login User

**Endpoint:**

```
POST /api/login
```

**Request Body:**

```json
{
  "email": "kherld@forum.com",
  "password": "frmpwd"
}
```

**Test with cURL:**

```
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "kherld@forum.com", "password": "frmpwd"}'
```

## Posts

### Create a Post (Authenticated)

**Endpoint:**

```
POST /api/posts
```

**Headers:**

```
Authorization: Bearer <token>
```

**Request Body:**

```json
{
  "title": "My First Post",
  "body": "This is the content of my first post."
}
```

**Test with cURL:**

```
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title": "My First Post", "body": "This is the content of my first post."}'
```

### Get All Posts

**Endpoint:**

```
GET /api/posts
```

**Test with cURL:**

```
curl -X GET http://localhost:8080/api/posts
```

### Like a Post (Authenticated)

**Endpoint:**

```
POST /api/posts/like
```

**Headers:**

```
Authorization: Bearer <token>
```

**Request Body:**

```json
{
  "post_id": 1
}
```

**Test with cURL:**

```
curl -X POST http://localhost:8080/api/posts/like \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"post_id": 1}'
```

## Comments

### Add a Comment to a Post (Authenticated)

**Endpoint:**

```
POST /api/posts/comment
```

**Headers:**

```
Authorization: Bearer <token>
```

**Request Body:**

```json
{
  "post_id": 1,
  "body": "This is my comment."
}
```

**Test with cURL:**

```
curl -X POST http://localhost:8080/api/posts/comment \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"post_id": 1, "body": "This is my comment."}'
```

## Sessions

### Validate Session Token

**Endpoint:**

```
GET /api/session/validate
```

**Headers:**

```
Authorization: Bearer <token>
```

**Test with cURL:**

```
curl -X GET http://localhost:8080/api/session/validate \
  -H "Authorization: Bearer <token>"
```

### **Testing Multiple Formats with `curl`**

Now you can test uploads with different file formats:

#### **Uploading an Image (PNG)**

```sh
curl -X POST http://localhost:8080/api/upload \
  -H "Content-Type: multipart/form-data" \
  -F "file=@path/to/image.png"
```

#### **Uploading a Video (MP4)**

```sh
curl -X POST http://localhost:8080/api/upload \
  -H "Content-Type: multipart/form-data" \
  -F "file=@path/to/video.mp4"
```

#### **Uploading an Audio File (MP3)**

```sh
curl -X POST http://localhost:8080/api/upload \
  -H "Content-Type: multipart/form-data" \
  -F "file=@path/to/audio.mp3"
```

#### **Uploading an Unsupported File (EXE)**

```sh
curl -X POST http://localhost:8080/api/upload \
  -H "Content-Type: multipart/form-data" \
  -F "file=@path/to/file.exe"
```

> **Expected Response:** `"unsupported file format: .exe"`
