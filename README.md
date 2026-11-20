# City Letter Counter

Enter a letter and see how many cities start with it. React frontend, Go backend.

The cities are the 50 most populated places from GeoNames, fetched with the query from the task. For example, `r` gives 1 (Rio de Janeiro) and `c` gives 3 (Cairo, Chongqing, Chengdu).

## Requirements

- Node.js 24
- Go 1.22+

## Run

```sh
npm run dev
```

Open <http://localhost:5173>. The Go API runs on port 3000, and Vite proxies `/api` to it.

## API

```
GET /api/cities?startsWith=r  ->  { "startsWith": "r", "count": 1, "cities": ["Rio de Janeiro"] }
GET /api/cities               ->  { "count": 50, "cities": [...] }
```

`startsWith` must be a single letter (case doesn't matter), otherwise you get a `400`.

The city list is cached in memory for 10 minutes. If GeoNames is down, the API keeps serving the last list for up to 24 hours and returns `502` after that.

Settings (environment variables):

- `GEONAMES_USERNAME`: defaults to `hsample`. It's a shared account and often hits its limit, so use your own if requests fail.
- `GEONAMES_ENDPOINT`: defaults to `https://secure.geonames.org/searchJSON`.
- `PORT`: defaults to `3000`.

## License

MIT, see [LICENSE](LICENSE).
