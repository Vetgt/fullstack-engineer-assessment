import React, { useEffect, useState } from 'react';
import { ActivityIndicator, Modal, Pressable, StyleSheet, Text, TextInput, View } from 'react-native';
import { Task, TaskStatus, UpdateTaskPayload } from '../types/task';

interface TaskEditModalProps {
  visible: boolean;
  task: Task | null;
  saving: boolean;
  onClose: () => void;
  onSave: (payload: UpdateTaskPayload) => Promise<void>;
}

const statuses: Array<{ label: string; value: TaskStatus }> = [
  { label: 'To Do', value: 'todo' },
  { label: 'In Progress', value: 'in_progress' },
  { label: 'Done', value: 'done' },
];

export function TaskEditModal({ visible, task, saving, onClose, onSave }: TaskEditModalProps) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [status, setStatus] = useState<TaskStatus>('todo');
  const [assignee, setAssignee] = useState('');

  useEffect(() => {
    if (!task) return;
    setTitle(task.title);
    setDescription(task.description);
    setStatus(task.status);
    setAssignee(task.assignee);
  }, [task]);

  const submit = async () => {
    if (!title.trim()) return;
    await onSave({
      title: title.trim(),
      description: description.trim(),
      status,
      assignee: assignee.trim(),
    });
  };

  return (
    <Modal visible={visible} transparent animationType="slide" onRequestClose={onClose}>
      <View style={styles.overlay}>
        <View style={styles.sheet}>
          <View style={styles.header}>
            <Text style={styles.heading}>Edit task</Text>
            <Pressable onPress={onClose} disabled={saving} style={styles.closeButton}>
              <Text style={styles.closeText}>Close</Text>
            </Pressable>
          </View>

          <Text style={styles.label}>Title</Text>
          <TextInput value={title} onChangeText={setTitle} style={styles.input} editable={!saving} />

          <Text style={styles.label}>Description</Text>
          <TextInput
            value={description}
            onChangeText={setDescription}
            multiline
            numberOfLines={4}
            textAlignVertical="top"
            style={[styles.input, styles.multiline]}
            editable={!saving}
          />

          <Text style={styles.label}>Assignee</Text>
          <TextInput value={assignee} onChangeText={setAssignee} style={styles.input} editable={!saving} />

          <Text style={styles.label}>Status</Text>
          <View style={styles.statusRow}>
            {statuses.map((option) => {
              const active = option.value === status;
              return (
                <Pressable
                  key={option.value}
                  disabled={saving}
                  onPress={() => setStatus(option.value)}
                  style={[styles.statusButton, active && styles.activeStatusButton]}
                >
                  <Text style={[styles.statusButtonText, active && styles.activeStatusButtonText]}>{option.label}</Text>
                </Pressable>
              );
            })}
          </View>

          <Pressable
            accessibilityRole="button"
            disabled={saving || !title.trim()}
            onPress={submit}
            style={[styles.saveButton, (saving || !title.trim()) && styles.disabledButton]}
          >
            {saving ? <ActivityIndicator color="#FFFFFF" /> : <Text style={styles.saveText}>Save changes</Text>}
          </Pressable>
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(15, 23, 42, 0.45)',
    justifyContent: 'flex-end',
  },
  sheet: {
    backgroundColor: '#FFFFFF',
    borderTopLeftRadius: 22,
    borderTopRightRadius: 22,
    padding: 20,
    gap: 8,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: 6,
  },
  heading: {
    fontSize: 20,
    fontWeight: '800',
    color: '#172033',
  },
  closeButton: {
    paddingVertical: 6,
    paddingHorizontal: 8,
  },
  closeText: {
    color: '#536174',
    fontWeight: '700',
  },
  label: {
    color: '#536174',
    fontWeight: '700',
    fontSize: 12,
    marginTop: 4,
  },
  input: {
    borderWidth: 1,
    borderColor: '#D7DDE7',
    borderRadius: 12,
    paddingHorizontal: 12,
    paddingVertical: 11,
    fontSize: 14,
    color: '#172033',
    backgroundColor: '#FFFFFF',
  },
  multiline: {
    minHeight: 90,
  },
  statusRow: {
    flexDirection: 'row',
    gap: 8,
  },
  statusButton: {
    flex: 1,
    borderWidth: 1,
    borderColor: '#D7DDE7',
    borderRadius: 10,
    paddingVertical: 10,
    alignItems: 'center',
  },
  activeStatusButton: {
    borderColor: '#1F2937',
    backgroundColor: '#1F2937',
  },
  statusButtonText: {
    fontSize: 12,
    color: '#536174',
    fontWeight: '700',
  },
  activeStatusButtonText: {
    color: '#FFFFFF',
  },
  saveButton: {
    marginTop: 8,
    backgroundColor: '#111827',
    borderRadius: 12,
    minHeight: 48,
    alignItems: 'center',
    justifyContent: 'center',
  },
  disabledButton: {
    opacity: 0.5,
  },
  saveText: {
    color: '#FFFFFF',
    fontWeight: '800',
  },
});
