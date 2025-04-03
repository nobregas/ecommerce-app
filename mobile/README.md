# E-commerce App

## How to Run

To set up and run the project, execute the following commands:

```
update your ipv4

npm install
npm run dev
npm run start
```

## Environment Variables

The application uses environment variables for configuration. These are defined in the following files:

### Method 1: Direct Configuration (Recommended for development)

For quick development, you can directly edit the environment values in:
- `mobile/utils/env.ts`

This file contains hardcoded values that will be used throughout the application.

### Method 2: Using .env File

For production or more flexible configuration:

1. Copy the `.env.sample` file to a new file named `.env`
2. Update the environment variables in `.env`:
   - `API_BASE_URL`: Your backend API URL (e.g., http://192.168.1.100:8080)
   - `GEMINI_API_KEY`: Your Gemini AI API key (required for AI support chat)
   - `GEMINI_API_URL`: The Gemini API endpoint

**Note:** When using Method 2, you'll need to rebuild the application for changes to take effect.

## Gemini API Key

If you don't have a Gemini API key, you can get one from:
- [Google AI Studio](https://ai.google.dev/) - Create an account and generate an API key

**Note:** Without a valid Gemini API key, the AI support chat will run in mock mode with predefined responses.

## OBS

To test the UI for this ecommerce app was used virtual Pixel_8

