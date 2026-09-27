import React from 'react';
import { Pressable, StyleSheet, Text, View } from 'react-native';
import { Task, TaskStatus } from '../types/task';

interface TaskCardProps {
  task: Task;
  onEdit: (task: Task) => void;
}

const statusMeta: Record<TaskStatus, { label: string; background: string; text: string }> = {
  todo: { label: 'To Do', background: '#EEF2FF', text: '#4338CA' },
  in_progress: { label: 'In Progress', background: '#FFF7ED', text: '#C2410C' },
  done: { label: 'Done', background: '#ECFDF5', text: '#047857' },
};

export function TaskCard({ task, onEdit }: TaskCardProps) {
  const meta = statusMeta[task.status];

  return (
    <View style={styles.card}>
      <View style={styles.headerRow}>
        <View style={styles.titleContainer}>
          <Text style={styles.title}>{task.title}</Text>
          <Text style={styles.assignee}>{task.assignee || 'Unassigned'}</Text>
        </View>
        <View style={[styles.status, { backgroundColor: meta.background }]}>
          <Text style={[styles.statusText, { color: meta.text }]}>{meta.label}</Text>
        </View>
      </View>

      {!!task.description && <Text style={styles.description}>{task.description}</Text>}

      <View style={styles.footerRow}>
        <Text style={styles.date}>{new Date(task.updated_at).toLocaleString()}</Text>
        <Pressable accessibilityRole="button" onPress={() => onEdit(task)} style={styles.editButton}>
          <Text style={styles.editText}>Edit</Text>
        </Pressable>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    borderRadius: 16,
    backgroundColor: '#FFFFFF',
    padding: 16,
    marginBottom: 12,
    borderWidth: 1,
    borderColor: '#E7EBF1',
  },
  headerRow: {
    flexDirection: 'row',
    gap: 10,
    alignItems: 'flex-start',
  },
  titleContainer: {
    flex: 1,
  },
  title: {
    color: '#172033',
    fontSize: 16,
    fontWeight: '700',
  },
  assignee: {
    color: '#738096',
    fontSize: 12,
    marginTop: 4,
  },
  status: {
    borderRadius: 999,
    paddingHorizontal: 10,
    paddingVertical: 6,
  },
  statusText: {
    fontSize: 11,
    fontWeight: '700',
  },
  description: {
    color: '#536174',
    fontSize: 13,
    lineHeight: 19,
    marginTop: 12,
  },
  footerRow: {
    marginTop: 14,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  date: {
    flex: 1,
    color: '#8A94A6',
    fontSize: 11,
    marginRight: 10,
  },
  editButton: {
    backgroundColor: '#F3F4F6',
    borderRadius: 10,
    paddingHorizontal: 13,
    paddingVertical: 9,
  },
  editText: {
    color: '#1F2937',
    fontWeight: '700',
    fontSize: 12,
  },
});
