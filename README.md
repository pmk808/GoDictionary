# GoDictionary

GoDictionary is a simple dictionary application that allows users to search for word definitions, pronunciations, idioms, and audio pronunciations. Built with React for the frontend and Go for the backend, this app provides a user-friendly interface for looking up words.

## Features
- Search for words and retrieve their definitions.
- Display pronunciations in standard and IPA formats.
- Provide audio pronunciations for words.
- List idioms related to the searched word.

## Technologies Used
- **Frontend**: React, CSS
- **Backend**: Go, Gorilla Mux
- **API**: Merriam-Webster Dictionary API

## Getting Started

### Prerequisites
- Go (version 1.16 or higher)
- Node.js (version 14 or higher)
- npm or yarn

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/godictionary.git
   cd godictionary
   ```

2. **Set up the backend**
   - Navigate to the backend directory:
     ```bash
     cd backend
     ```
   - Install dependencies:
     ```bash
     go mod tidy
     ```
   - Run the backend server:
     ```bash
     go run main.go
     ```
   - The server will start on `http://localhost:8080`.

3. **Set up the frontend**
   - Navigate to the frontend directory:
     ```bash
     cd frontend
     ```
   - Install dependencies:
     ```bash
     npm install
     ```
   - Start the frontend application:
     ```bash
     npm start
     ```
   - The application will be available at `http://localhost:3000`.

## API Endpoints
### `GET /word`
- **Description**: Fetches word information based on the provided query parameter.
- **Query Parameters**:
  - `text`: The word to look up.
- **Response**: Returns a JSON object containing the word's definitions, pronunciations, idioms, and audio links.

### Example Request

- GET http://localhost:8080/word?text=example


### Example Response

- json
```{
"text": "example",
"pronunciations": ["ig-zam-puhl"],
"ipaPronunciation": "/ɪɡˈzæmpəl/",
"audioPronunciations": ["audio1", "audio2"],
"definitions": [
{
"partOfSpeech": "noun",
"senses": ["a representative form or pattern", "something to be imitated"]
}
],
"idioms": [
{
"phrase": "for example",
"senses": ["used to introduce a case or instance"]
}
]
}
```

## Contributing
Contributions are welcome! Please feel free to submit a pull request or open an issue for any suggestions or improvements.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments
- [Merriam-Webster Dictionary API](https://dictionaryapi.com/)
- [Gorilla Mux](https://github.com/gorilla/mux)
- [React](https://reactjs.org/)



