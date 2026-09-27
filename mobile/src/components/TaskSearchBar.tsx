import React from 'react';
import { StyleSheet, TextInput, View } from 'react-native';

interface TaskSearchBarProps {
  value: string;
  onChangeText: (value: string) => void;
}

export function TaskSearchBar({ value, onChangeText }: TaskSearchBarProps) {
  return (
    <View style={styles.container}>
      <TextInput
        accessibilityLabel="Search tasks"
        value={value}
        onChangeText={onChangeText}
        placeholder="Search title or description"
        placeholderTextColor="#8A94A6"
        autoCorrect={false}
        returnKeyType="search"
        style={styles.input}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginBottom: 12,
  },
  input: {
    borderWidth: 1,
    borderColor: '#D7DDE7',
    borderRadius: 12,
    backgroundColor: '#FFFFFF',
    paddingHorizontal: 14,
    paddingVertical: 12,
    fontSize: 15,
    color: '#1D2633',
  },
});
