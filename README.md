# Interview Challenge for Joe Edwards

## Set Config
conf.env contains configurable values.  Replace existing values with preferences.

## To Run the Server
```bash
go build -o weather_service .

go run .
```

## Testing the Service for Kansas City
Once the server is running, you can use this CURL command to test weather in Kansas City.  Choosing points for other cities has been unreliable.  Ensure that the same port used in the CURL command is the same port as defined in conf.env.


```bash
# Test with curl
curl -X POST http://localhost:8080/weather \
  -H "Content-Type: application/json" \
  -d '{"location": {"latitude": "39.0997", "longitude": "-94.5786"}}'
```


## Data Models

### Request Body
```json
{
  "location": {
    "latitude": "39.0997",
    "longitude": "-94.5786"
  }
}
```

### Response Body
```json
{
  "status": "success",
  "forecast": "Tonight: Clear. Temperature: 45°F. Wind: 5 mph NW",
  "characterization": "Cold"
}
```

### Error Response Body
```json
{
  "status": "failure",
  "error_code": "Failed to reach NWS points API: connection timeout"
}
```

### Shortcuts
- Leveraged AI to build boilerplate for config, JSON model conversion, and basic server organization.
- Config values for maximums are in Fahrenheit only.