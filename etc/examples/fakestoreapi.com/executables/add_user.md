* username:"mor_2314"
* password:"83r5^_"

```http
POST https://fakestoreapi.com/users HTTP/1.1
Content-Type: application/json
Content-Length: 277

{"email":"John@gmail.com","username":"{{username}}","password":"{{password}}","name":{"firstname":"John","lastname":"Doe"},"address":{"city":"kilcoole","street":"7835 new road","number":3,"zipcode":"12926-3874","geolocation":{"lat":"-37.3159","long":"81.1496"}},"phone":"1-570-236-7033"}
```

* new_user_id:`Body.id`