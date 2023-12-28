# TXT to SRT Converter

This is a simple web application written in Go using the Gin framework that allows users to convert text files into SubRip (SRT) subtitle format each line has a duration of 4 second by default.

## Features

- **User-Friendly Interface:** The web application provides a clean and intuitive interface for users to interact with.
- **File Upload:** Users can upload a text file using the provided file input.
- **Automatic Conversion:** The uploaded text file is automatically converted to the SubRip (SRT) format.
## Getting Started

### Prerequisites

- Go installed on your machine.

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/your-username/subtitle-converter.git
    ```

2. Navigate to the project directory:

    ```bash
    cd subtitle-converter
    ```

3. Run the application:

    ```bash
    go run main.go
    ```

4. Open your web browser and go to [http://localhost:8080](http://localhost:8080) to access the application.

5. Try it on the follow link [https://txt-to-srt.onrender.com](https://txt-to-srt.onrender.com)

## Usage

1. Open the application in your web browser.
2. Click on the "Choose File" button and select the text file you want to convert.
3. Click the "Convert" button to initiate the conversion process.

## License

This project is licensed under the [MIT License](LICENSE).