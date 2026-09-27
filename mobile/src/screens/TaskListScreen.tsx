import React, { useCallback, useEffect, useMemo, useState } from 'react';
import {
  ActivityIndicator,
  Alert,
  FlatList,
  RefreshControl,
  SafeAreaView,
  StyleSheet,
  Text,
  View,
} from 'react-native';

import { TaskCard } from '../components/TaskCard';
import { TaskEditModal } from '../components/TaskEditModal';
import { TaskSearchBar } from '../components/TaskSearchBar';
import { StatusFilter } from '../components/StatusFilter';
import { getTasks, updateTask } from '../services/api';
import { Task, TaskFilters, TaskStatus, TaskListResponse, UpdateTaskPayload } from '../types/task';

const DEFAULT_FILTERS: TaskFilters = {
  page: 1,
  limit: 5,
  keyword: '',
  status: '',
  assignee: '',
  sort: 'created_at_desc',
};

export function TaskListScreen() {
  const [filters, setFilters] = useState<TaskFilters>(DEFAULT_FILTERS);
  const [data, setData] = useState<TaskListResponse>({ items: [], page: 1, limit: 5, total: 0, total_pages: 0 });
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [saving, setSaving] = useState(false);
  const [selectedTask, setSelectedTask] = useState<Task | null>(null);

  const loadTasks = useCallback(async () => {
    try {
      setLoading(true);
      const next = await getTasks(filters);
      setData(next);
    } catch (error) {
      Alert.alert('Could not load tasks', error instanceof Error ? error.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    const timer = setTimeout(loadTasks, filters.keyword ? 300 : 0);
    return () => clearTimeout(timer);
  }, [loadTasks, filters.keyword]);

  const handleRefresh = async () => {
    setRefreshing(true);
    await loadTasks();
    setRefreshing(false);
  };

  const handleSearch = (keyword: string) => {
    setFilters((current) => ({ ...current, keyword, page: 1 }));
  };

  const handleStatus = (status: '' | TaskStatus) => {
    setFilters((current) => ({ ...current, status, page: 1 }));
  };

  const changePage = (delta: number) => {
    setFilters((current) => ({
      ...current,
      page: Math.max(1, Math.min(data.total_pages, current.page + delta)),
    }));
  };

  const saveTask = async (payload: UpdateTaskPayload) => {
    if (!selectedTask) return;

    try {
      setSaving(true);
      await updateTask(selectedTask.id, payload);
      setSelectedTask(null);
      await loadTasks();
      Alert.alert('Task updated', 'Changes were saved successfully.');
    } catch (error) {
      Alert.alert('Update failed', error instanceof Error ? error.message : 'Unknown error');
    } finally {
      setSaving(false);
    }
  };

  const footerText = useMemo(() => {
    if (data.total === 0) return 'No tasks found';
    return `Page ${data.page} of ${data.total_pages} · ${data.total} total`;
  }, [data]);

  return (
    <SafeAreaView style={styles.safeArea}>
      <FlatList
        data={data.items}
        keyExtractor={(item) => String(item.id)}
        contentContainerStyle={styles.listContent}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={handleRefresh} />}
        ListHeaderComponent={
          <View>
            <Text style={styles.title}>Task Management</Text>
            <Text style={styles.subtitle}>Search, filter, paginate and edit tasks from the Go API.</Text>
            <TaskSearchBar value={filters.keyword} onChangeText={handleSearch} />
            <StatusFilter value={filters.status} onChange={handleStatus} />
          </View>
        }
        renderItem={({ item }) => <TaskCard task={item} onEdit={setSelectedTask} />}
        ListEmptyComponent={
          loading ? (
            <View style={styles.emptyState}>
              <ActivityIndicator size="large" />
              <Text style={styles.emptyText}>Loading tasks...</Text>
            </View>
          ) : (
            <View style={styles.emptyState}>
              <Text style={styles.emptyText}>No tasks match the current filters.</Text>
            </View>
          )
        }
        ListFooterComponent={
          <View style={styles.footer}>
            <Text style={styles.footerText}>{footerText}</Text>
            <View style={styles.paginationRow}>
              <Text
                accessibilityRole="button"
                onPress={() => data.page > 1 && changePage(-1)}
                style={[styles.pageButton, data.page <= 1 && styles.disabledPageButton]}
              >
                Previous
              </Text>
              <Text
                accessibilityRole="button"
                onPress={() => data.page < data.total_pages && changePage(1)}
                style={[styles.pageButton, data.page >= data.total_pages && styles.disabledPageButton]}
              >
                Next
              </Text>
            </View>
          </View>
        }
      />

      <TaskEditModal
        visible={Boolean(selectedTask)}
        task={selectedTask}
        saving={saving}
        onClose={() => !saving && setSelectedTask(null)}
        onSave={saveTask}
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
    backgroundColor: '#F6F8FB',
  },
  listContent: {
    paddingHorizontal: 16,
    paddingTop: 12,
    paddingBottom: 32,
  },
  title: {
    fontSize: 28,
    fontWeight: '900',
    color: '#172033',
  },
  subtitle: {
    marginTop: 6,
    marginBottom: 18,
    color: '#738096',
    lineHeight: 19,
    fontSize: 13,
  },
  emptyState: {
    minHeight: 220,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  emptyText: {
    marginTop: 12,
    textAlign: 'center',
    color: '#738096',
  },
  footer: {
    paddingTop: 10,
    paddingBottom: 20,
  },
  footerText: {
    textAlign: 'center',
    color: '#738096',
    fontSize: 12,
  },
  paginationRow: {
    flexDirection: 'row',
    justifyContent: 'center',
    gap: 28,
    marginTop: 14,
  },
  pageButton: {
    color: '#111827',
    fontWeight: '800',
  },
  disabledPageButton: {
    color: '#B7BFCC',
  },
});
