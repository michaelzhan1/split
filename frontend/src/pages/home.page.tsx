import { useMutation } from '@tanstack/react-query';

import type { AxiosError } from 'axios';
import { useState } from 'react';
import { useNavigate } from 'react-router';
import { createGroup } from 'src/services/group.service';

export function Home() {
  const navigate = useNavigate();

  const [groupId, setGroupId] = useState<string>('');
  const [groupName, setGroupName] = useState<string>('');

  const { mutate: createGroupMutate, isPending: isPendingCreateGroup } =
    useMutation<{ id: number }, AxiosError, { name: string }>({
      mutationFn: (variables: { name: string }) => {
        return createGroup(variables.name);
      },
    });
  const onCreateGroup = (name: string) =>
    createGroupMutate(
      { name },
      {
        onSuccess: (data) => {
          navigate(`/groups/${data.id}`);
        },
        onError: (error) => {
          console.error('Error creating group:', error);
          alert('Failed to create group. Please try again.');
        },
      },
    );

  return (
    <>
      <h1>Split</h1>
      <h3>Find a group</h3>
      <form>
        <label htmlFor='group-code'>Group Code</label>
        <input
          name='group-code'
          type='number'
          placeholder='Group Code'
          required
          value={groupId}
          onChange={(e) => setGroupId(e.target.value)}
        />
        <button
          type='submit'
          onClick={(e) => {
            e.preventDefault();
            navigate(`/groups/${groupId}`);
          }}
        >
          Go
        </button>
      </form>

      <h3>Create a new group</h3>
      <form>
        <label htmlFor='group-name'>Group Name</label>
        <input
          name='group-name'
          type='text'
          placeholder='Group Name'
          required
          value={groupName}
          onChange={(e) => setGroupName(e.target.value)}
        />
        <button
          type='submit'
          onClick={(e) => {
            e.preventDefault();
            onCreateGroup(groupName);
          }}
          disabled={isPendingCreateGroup}
        >
          Create
        </button>
      </form>
    </>
  );
}
