# รายชื่อผลิตภัณฑ์ทั้งหมด

ทดสอบ API สำหรับลิสต์รายชื่อผลิตภัณฑ์ทั้งหมด โดยจะต้องทำการล็อกอินเข้าใช้งานเพื่อเอา token ก่อน และส่ง token ไปใน header เพื่อเรียกดูรายชื่อผลิตภัณฑ์

* [ลงชื่อเข้าใช้ระบบ](../executables/auth.md)

```http
GET https://fakestoreapi.com/products HTTP/1.1
Authorization: Bearer {{access_token}}
```

* `Body.0.id`==`1`