# URL Shortener using Gin

This service will take a long URL and generates a short URL. It also provides the functionality to retrieve the original URL when given the short URL.

### Endpoints:
- POST `/shorten`: This endpoint accepts a JSON request body in the format `{"url": "<long_url>"}` and responds with a JSON object containing the shortened URL, such as `{"short_url": "<short_url>"}`.
- GET `/{short_url}`: This endpoint takes the short URL as a path parameter and redirects the user to the original URL.

### Features:
- Implemented the service in a single `main.go` file.
- Used a `hashing method` to generate  short urls.
- Used `maps` for storing short URLs and their original URLs.
- Added simple URL validation to ensure the provided long URL is valid by checking if it uses HTTPS.
- The service is thread-safe, ensuring safe concurrent access.
### Running the application :
```bash
go run main.go
```
### Testing :
Test cases have been added to `main_test.go`. You can run these tests using the following command:

```bash
go test
```

