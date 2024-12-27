### How to use it:

1. **Create `.env` File**: 
   In the root directory of your project, create a file named `.env` and add the environment variables as shown above. This file contains your database and JWT configuration settings.
   
2. **Install Dependencies**:
   The `go mod tidy` command ensures that your Go modules and dependencies are correctly installed.

3. **Run the Application**:
   Use `go run main.go` to run your application locally. This will start your API server on `http://localhost:8080`.

4. **Access Swagger UI**:
   You can open the Swagger UI at `http://localhost:8080/swagger-ui` to view the API documentation and interact with the API.
