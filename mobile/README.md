# Fullstack Engineer Assessment - React Native Frontend

React Native + TypeScript client for the task management assessment.

## Implemented requirements

- Search input
- Status filter
- Pagination
- Edit modal using `PUT /api/tasks/:id`
- Loading state
- Pull-to-refresh after a refresh action
- Error state via native alerts
- At least one component test for the search input
- TypeScript types for API contracts

## Requirements

- Node.js 22.13+ for Expo SDK 57
- npm
- Expo CLI through `npx`
- iOS Simulator, Android Emulator, or a physical device

## Install

```bash
npm install
```

Copy `.env.example` to `.env` when you need a custom API URL. Expo automatically exposes variables prefixed with `EXPO_PUBLIC_` to application code.

### iOS Simulator

Use:

```text
EXPO_PUBLIC_API_URL=http://localhost:8080/api
```

Then:

```bash
npm start
```

Press `i`.

### Android Emulator

Use:

```text
EXPO_PUBLIC_API_URL=http://10.0.2.2:8080/api
```

Then:

```bash
npm start
```

Press `a`.

### Physical phone

The phone must reach the computer running Go. Do not use `localhost` from the phone.

Find your Mac LAN IP, for example:

```bash
ipconfig getifaddr en0
```

Set:

```text
EXPO_PUBLIC_API_URL=http://192.168.1.20:8080/api
```

Replace the IP with your actual LAN address, then restart Expo.

## Test

```bash
npm test
```

## Type check

```bash
npm run typecheck
```
