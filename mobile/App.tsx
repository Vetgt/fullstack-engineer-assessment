import React from 'react';
import { StatusBar } from 'expo-status-bar';
import { TaskListScreen } from './src/screens/TaskListScreen';

export default function App() {
  return (
    <>
      <StatusBar style="dark" />
      <TaskListScreen />
    </>
  );
}
