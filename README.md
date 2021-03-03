# userdb-go
## Fictitious Http server to maintain user data

## Syntax:

#### **URL** : /create/
#### **POST** data:
One or several objects consisting of:
* Name -- a unique name
*	Passwd -- password (non-encrypted at this point)
*	Email -- user's email (optional)
* Addr  -- user's street address (optional)

**Response**: a list of created UIDs and/or a error message

**Examples**:

*Single mode*
```
> curl --request POST localhost:8080/create/ -d "{\"name\":\"ilya\", \"passwd\":\"asdfgh\", \"email\":\"Bububu@hhh\"}"
{"Uid":1001}
```

*Batch mode*  (one of the names is repeated, so its request and ones after it are ignored)
```
> curl --request POST localhost:8080/create/ -d "{\"name\":\"ilya1\", \"passwd\":\"a\"}{\"name\":\"ilya2\", \"passwd\":\"a\"}{\"name\":\"ilya3\", \"passwd\":\"a\"}{\"name\":\"ilya4\", \"passwd\":\"a\"}{\"name\":\"ilya4\", \"passwd\":\"a\"}{\"name\":\"ilya5\", \"passwd\":\"a\"}"
[{"Uid":1002},{"Uid":1003},{"Uid":1004},{"Uid":1005},{"Uid":-1}]
{"Error":"Name already reserved"}
```

--------------------------------------------------------------------------------------------------------

#### **URL** : /update/_**UID**_
#### **POST** data:
A single object consisting of one or both of:
*	Email -- update for user's email
* Addr  -- update for user's street address

**Response**: Status OK or a Error message

**Examples** :

```
> curl --request POST localhost:8080/update/1004 -d "{\"email\":\"ilya@hhh\", \"addr\": \"123 Bway\" }"
{"Status":"OK"}

> curl --request POST localhost:8080/update/9999 -d "{\"email\":\"ilya@hhh\", \"addr\": \"123 Bway\" }"
{"Status":"ERROR: UID not found"}
```

--------------------------------------------------------------------------------------------------------

#### **URL** : /get/_**UID**_

**Response**: User record object or a Error message

**Examples** :

```
> curl --request GET localhost:8080/get/1004
{"Name":"ilya3","Email":"ilya@hhh","Addr":"123 Bway"}

> curl --request GET localhost:8080/get/9999
{"Error":"UID not found"}
```
