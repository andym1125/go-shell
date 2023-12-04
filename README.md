# Project 2: Shell Builtins
Author: Andy McDowall
## Testing

For testing and code coverage we use the following command:
```bash
go test ./... -v -short -p 1 -cover
```

## Example run

```bash
(base) andy@Andys-MacBook-Pro-3 Project2 % go run main.go
/Users/andy/Documents/workshop/CSCE4600/Project2 [andy] $ echo hello
hello
/Users/andy/Documents/workshop/CSCE4600/Project2 [andy] $ echo -n hello
hello/Users/andy/Documents/workshop/CSCE4600/Project2 [andy] $ repeat 3 echo hello
hello
hello
hello
/Users/andy/Documents/workshop/CSCE4600/Project2 [andy] $ pwd 
/Users/andy/Documents/workshop/CSCE4600/Project2
/Users/andy/Documents/workshop/CSCE4600/Project2 [andy] $ pwd -P
/Users/andy/Documents/workshop/CSCE4600/Project2
/Users/andy/Documents/workshop/CSCE4600/Project2 [andy] $ history 3
66 repeat 3 echo hello
67 pwd 
68 pwd -P
/Users/andy/Documents/workshop/CSCE4600/Project2 [andy] $ history 3 -rh
history 3
pwd -P
pwd 
/Users/andy/Documents/workshop/CSCE4600/Project2 [andy] $ export hello=world
/Users/andy/Documents/workshop/CSCE4600/Project2 [andy] $ env
<output censored>
HISTFILE=./.shell/history.txt
hello=world
/Users/andy/Documents/workshop/CSCE4600/Project2 [andy] $ exit
exiting gracefully...
(base) andy@Andys-MacBook-Pro-3 Project2 % 
```

## Source

Thank you to:
https://stackoverflow.com/questions/10516662/how-to-measure-test-coverage-in-go