# Clog
This is an example of using [Clarifai](https://clarifai.com/)'s image recognition API with Golang. The application uploads images to Cloudinary and the passes the link to Clarify for analyzes.

Note that only small(ish) files can be uploaded and analyzed.

To use the API you need to:
- Have an account with Clarifai
- Create a project at Clarifai
- Generate an API key.
- Set the key as and environment variable with the name CLARIFAI_TEST_APP

## Todo
Resize large images before uploading.

