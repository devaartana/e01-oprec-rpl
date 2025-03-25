# Link Shortener API
Implementasi ini membuat sebuah API yang dapat digunakan untuk membuat sebuah webapp link shortener dimana terdapat user yang dapat memiliki/membuat sebuah link shortener, berikut fiturnya
- Register akun
- Login akun
- CRUD link
- Update expired date link
- Redirect 

## Database
MongoDB saya pilih sebagai database untuk API ini karena memungkinkan penyimpanan data yang fleksibel dengan mengembed data links dalam satu dokumen per pengguna dan data yang disimpan tidak memiliki relasi yang banyak sehingga membuat lebih mudah dikelola.

Contoh dokumen
```
{
  "_id": ObjectId("64b5f7e8d9f3a1c9a7e3b123"),
  "username": "kadekarya",
  "password": "hashedpassword123",
  "email": "deva@example.com",
  "created_at": ISODate(""),
  "links": [
    {
      "slug": "xyz123",
      "original_url": "https://example.com/long-url-example",
      "created_at": ISODate(""),
      "expired_date": ISODate(""),

    },
    {
      "slug": "abc456",
      "original_url": "https://example2.com/another-url",
      "created_at": ISODate(""),
      "expired_date": ISODate(""),
    }
  ]
}
```

Dengan diagram

![image](https://github.com/user-attachments/assets/e6537f7c-a1d0-4b4c-bed9-911d91e4e2c0)

Dalam desain ini, jika dilihat dari relasinya user memiliki relasi one-to-many ke links tetapi dalam implementasi ini data links diembed langsung ke dokumen user sebagai sebuah array.

## Dokumentasi API
- Register akun [POST]
  - Endpoint: localhost:8000/api/register
  - Request:
    ```
     {
      "Content-Type": "application/json",
      "Body": {
          "username": "Arya",
          "password": "password123",
          "email": "deva@gmail.com"
      }
    }
    ```
  - Response Success(201)
    ```
      message 
    ```
- Login akun [POST]
  - Endpoint: localhost:8000/api/login
  - Request:
    ```
    {
      "Content-Type": "application/json",
      "Body": {
            "email": "deva@gmail.com"
            "password": "password123",
        }
    }
    ```
  - Response Success(200)
    ```
      "token": abcd
    ```
- Create link [POST]
  - Endpoint: localhost:8000/api/links
  - Request:
    ```
    {
      "Content-Type": "application/json",
      "Authorization": "Bearer abcd"
      "Body": {
          "slug": "nice-king",
          "original_url": "www.youtube.com"
      }
    }
    ```
  - Response Success(201)
    ```
      message
    ```
- Read link [GET]
  - Endpoint: localhost:8000/api/links
  - Request:
    ```
    {
      "Authorization": "Bearer abcd"
    }
    ```
  - Response Success(200)
    ```
      message
    ```
- Update link [PUT]
  - Endpoint: localhost:8000/api/links
  - Request:
     ```
    {
      "Content-Type": "application/json",
      "Authorization": "Bearer abcd"
      "Body": {
          "slug": "nice-king",
          "original_url": "www.youtube.com"
      }
    }
    ```
  - Response Success(200)
    ```
      message
    ```
- Delete link [DELETE]
  - Endpoint: localhost:8000/api/links/{slug}
  - Request:
    ```
    {
      "Authorization": "Bearer abcd"
    }
    ```
  - Response Success(200)
    ```
      message
    ```
- Refresh expired date [GET]
  - Endpoint: localhost:8000/api/links/{slug}
  - Request:
    ```
    {
      "Authorization": "Bearer abcd"
    }
    ```
  - Response Success(200)
    ```
      message
    ```
- Redirect [GET]
  - Endpoint: localhost:8000/{slug}
