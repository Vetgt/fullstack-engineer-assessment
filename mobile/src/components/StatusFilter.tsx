import React from 'react';
import { ScrollView, StyleSheet, Text, Pressable } from 'react-native';
import { TaskStatus } from '../types/task';

const options: Array<{ label: string; value: '' | TaskStatus }> = [
  { label: 'All', value: '' },
  { label: 'To Do', value: 'todo' },
  { label: 'In Progress', value: 'in_progress' },
  { label: 'Done', value: 'done' },
];

interface StatusFilterProps {
  value: '' | TaskStatus;
  onChange: (value: '' | TaskStatus) => void;
}

export function StatusFilter({ value, onChange }: StatusFilterProps) {
  return (
    <ScrollView horizontal showsHorizontalScrollIndicator={false} contentContainerStyle={styles.content}>
      {options.map((option) => {
        const active = option.value === value;
        return (
          <Pressable
            key={option.label}
            accessibilityRole="button"
            accessibilityState={{ selected: active }}
            onPress={() => onChange(option.value)}
            style={[styles.chip, active && styles.activeChip]}
          >
            <Text style={[styles.text, active && styles.activeText]}>{option.label}</Text>
          </Pressable>
        );
      })}
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  content: {
    gap: 8,
    paddingBottom: 12,
  },
  chip: {
    borderRadius: 999,
    borderWidth: 1,
    borderColor: '#D7DDE7',
    backgroundColor: '#FFFFFF',
    paddingHorizontal: 14,
    paddingVertical: 9,
  },
  activeChip: {
    backgroundColor: '#1F2937',
    borderColor: '#1F2937',
  },
  text: {
    fontSize: 13,
    color: '#536174',
    fontWeight: '600',
  },
  activeText: {
    color: '#FFFFFF',
  },
});
