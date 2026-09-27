import React from 'react';
import { fireEvent, render } from '@testing-library/react-native';
import { TaskSearchBar } from '../src/components/TaskSearchBar';

describe('TaskSearchBar', () => {
  it('calls onChangeText when the user enters a search keyword', async () => {
    const onChangeText = jest.fn();

    const { getByPlaceholderText } = await render(
      <TaskSearchBar value="" onChangeText={onChangeText} />
    );

    const input = getByPlaceholderText('Search title or description');

    fireEvent.changeText(input, 'login');

    expect(onChangeText).toHaveBeenCalledWith('login');
  });
});