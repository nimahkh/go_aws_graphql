#### Config AWS 
create a file name `.aws/credentials` inside your root folder. 
``` mkdir ~/.aws && nano ~/.aws/credentials ```
and put these lines in your file as default settings of aws credentials

```
[default]
aws_access_key_id = <ACCESS KEY>
aws_secret_access_key = <ECRET ACCESS KEY>
```

### Install modules
``` go install ```

### Run the project

``` go run cmd/start.go ```

Now open `localhost:8085/users` on your browser and pass the GraphQl queries in URL as below:
```
http://localhost:8085/users?query={list{name}}
```