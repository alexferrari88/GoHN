# gohn — Hacker News API Wrapper for Go

<div align="center">
<a href="https://pkg.go.dev/github.com/alexferrari88/gohn"><img src="https://pkg.go.dev/badge/github.com/alexferrari88/gohn.svg" alt="Go Reference"></a>
</div>
<div align="center">
<img src="img/logo_1.svg" width="300" style="margin: 0 auto;" />
</div>

gohn is a tiny wrapper for the [Hacker News API](https://github.com/HackerNews/API) for Go.

It facilitates the use of the API by providing a simple interface to the API endpoints.

## Features 🚀

- Get the top stories
- Get the new stories
- Get the best stories
- Get the ask stories
- Get the show stories
- Get the job stories
- Retrieve all comments for a story using goroutines to speed up the process
- Apply filters to retrieved items (stories, comments)

## Usage 💻

Refer to the [GoDoc](https://pkg.go.dev/github.com/alexferrari88/gohn) for the full API reference.

### Example

Refer to [example/main.go](example/main.go) for a full example on how to retrieve the top stories and the all the comments for the first one.

## Contributing 🤝🏼

Feel free to fork this repo and create a PR. I will review them and merge if ok.
The above todos can be a very good place to start.

## License 📝

[MIT](https://choosealicense.com/licenses/mit/)
