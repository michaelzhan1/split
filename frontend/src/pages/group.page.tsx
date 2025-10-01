import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router';

import { skipToken, useMutation, useQuery } from '@tanstack/react-query';
import type { AxiosError } from 'axios';

import { AddUserModal } from 'src/components/add-user-modal.component';
import { ConfirmationModal } from 'src/components/confirmation-modal.component';
import { PatchGroupModal } from 'src/components/patch-group-modal.component';
import { PatchUserModal } from 'src/components/patch-user-modal.component';
import {
  deleteGroup,
  getGroupById,
  patchGroup,
} from 'src/services/group.service';
import {
  addUserToGroup,
  deleteUser,
  getUsersByGroupId,
  patchUser,
} from 'src/services/user.service';
import type { Group, User } from 'src/types/common.type';

export function Group() {
  const { groupId = '' } = useParams();
  const navigate = useNavigate();
  const [patchGroupModalOpen, setPatchGroupModalOpen] =
    useState<boolean>(false);
  const [deleteGroupModalOpen, setDeleteGroupModalOpen] =
    useState<boolean>(false);
  const [addUserModalOpen, setAddUserModalOpen] = useState<boolean>(false);
  const [patchUserModalOpen, setPatchUserModalOpen] =
    useState<boolean>(false);
  const [deleteUserModalOpen, setDeleteUserModalOpen] =
    useState<boolean>(false);

  const [selectedUser, setSelectedUser] = useState<User | null>(null);

  // group info
  const {
    data: group = null,
    isFetching: isFetchingGroup,
    refetch: refetchGroup,
    error: groupError,
  } = useQuery<Group, AxiosError>({
    queryKey: ['group', groupId],
    queryFn: () => getGroupById(Number(groupId)),
  });

  useEffect(() => {
    if (groupError) {
      console.error('Error fetching group:', groupError);
      alert(`Failed to fetch group: ${groupError.message}`);
      navigate('/');
    }
  }, [groupError, navigate]);

  // patch a group
  const { mutate: patchGroupMutate, isPending: isPendingPatchGroup } =
    useMutation<void, AxiosError, { name: string }>({
      mutationFn: (variables: { name: string }) => {
        return patchGroup(Number(groupId), variables.name);
      },
    });
  const onPatchGroup = (name: string) =>
    patchGroupMutate(
      { name },
      {
        onSuccess: () => {
          refetchGroup();
          setPatchGroupModalOpen(false);
        },
        onError: (error) => {
          console.error('Error updating group:', error);
          alert('Failed to update group. Please try again');
        },
      },
    );

  // delete group
  const { mutate: deleteGroupMutate, isPending: isPendingDeleteGroup } =
    useMutation<void, AxiosError>({
      mutationFn: () => {
        return deleteGroup(Number(groupId));
      },
    });
  const onDeleteGroup = () =>
    deleteGroupMutate(undefined, {
      onSuccess: () => {
        navigate('/');
      },
      onError: (error) => {
        console.error('Error deleting group:', error);
        alert('Failed to delete group. Please try again');
      },
    });

  // user info
  const {
    data: users = [],
    isFetching: isFetchingUsers,
    refetch: refetchUsers,
    error: usersError,
  } = useQuery<User[], AxiosError>({
    queryKey: ['users', groupId],
    queryFn: group ? () => getUsersByGroupId(group.id) : skipToken,
  });

  useEffect(() => {
    if (usersError) {
      console.error('Error fetching users:', usersError);
      alert('Failed to fetch users. Please try again.');
    }
  }, [usersError]);

  // add a user
  const { mutate: addUserMutate, isPending: isPendingAddUser } =
    useMutation<{ id: number }, AxiosError, { name: string }>({
      mutationFn: (variables: { name: string }) => {
        return addUserToGroup(Number(groupId), variables.name);
      },
    });
  const onAddUser = (name: string) =>
    addUserMutate(
      { name },
      {
        onSuccess: () => {
          refetchUsers();
          setAddUserModalOpen(false);
        },
        onError: (error) => {
          console.error('Error adding user:', error);
          alert('Failed to add user. Please try again');
        },
      },
    );

  // patch a user
  const { mutate: patchUserMutate, isPending: isPendingPatchUser } =
    useMutation<void, AxiosError, { userId: number; name: string }>({
      mutationFn: (variables: { userId: number; name: string }) => {
        return patchUser(Number(groupId), variables.userId, variables.name);
      },
    });
  const onPatchUser = (userId: number, name: string) =>
    patchUserMutate(
      { userId: userId, name },
      {
        onSuccess: () => {
          refetchUsers();
          setPatchUserModalOpen(false);
          setSelectedUser(null);
        },
        onError: (error) => {
          console.error('Error updating user:', error);
          alert('Failed to update user. Please try again');
        },
      },
    );

  // delete a user
  const { mutate: deleteUserMutate, isPending: isPendingDeleteUser } =
    useMutation<void, AxiosError, { userId: number }>({
      mutationFn: (variables: { userId: number }) => {
        return deleteUser(Number(groupId), variables.userId);
      },
    });
  const onDeleteUser = (userId: number) =>
    deleteUserMutate(
      { userId: userId },
      {
        onSuccess: () => {
          refetchUsers();
          setPatchUserModalOpen(false);
          setSelectedUser(null);
        },
        onError: (error) => {
          console.error('Error deleting user:', error);
          alert('Failed to delete user. Please try again');
        },
      },
    );

  const isLoading =
    isFetchingGroup ||
    isPendingPatchGroup ||
    isPendingDeleteGroup ||
    isFetchingUsers ||
    isPendingAddUser ||
    isPendingPatchUser ||
    isPendingDeleteUser;

  return !group || isLoading ? (
    <div>Loading...</div>
  ) : (
    <>
      <PatchGroupModal
        isOpen={patchGroupModalOpen}
        onClose={() => setPatchGroupModalOpen(false)}
        onSubmit={(name: string) => onPatchGroup(name)}
        initialName={group.name}
      />
      <ConfirmationModal
        isOpen={deleteGroupModalOpen}
        onClose={() => setDeleteGroupModalOpen(false)}
        title='Delete Group'
        content='Are you sure you want to delete this group? This action cannot be undone.'
        onSubmit={() => {
          onDeleteGroup();
          setDeleteGroupModalOpen(false);
        }}
      />
      <AddUserModal
        isOpen={addUserModalOpen}
        onClose={() => setAddUserModalOpen(false)}
        onSubmit={(name: string) => onAddUser(name)}
      />

      {selectedUser && (
        <>
          <PatchUserModal
            isOpen={patchUserModalOpen}
            onClose={() => setPatchUserModalOpen(false)}
            onSubmit={(userId: number, name: string) =>
              onPatchUser(userId, name)
            }
            user={selectedUser}
          />
          <ConfirmationModal
            isOpen={deleteUserModalOpen}
            onClose={() => setDeleteUserModalOpen(false)}
            title='Delete User'
            content={`Are you sure you want to delete user "${selectedUser.name}"? This action cannot be undone.`}
            onSubmit={() => {
              onDeleteUser(selectedUser.id);
              setDeleteUserModalOpen(false);
            }}
          />
        </>
      )}
      <div>
        <h1>Group: {group.name}</h1>
        <button onClick={() => setPatchGroupModalOpen(true)}>
          Edit group name
        </button>
        <button onClick={() => setDeleteGroupModalOpen(true)}>
          Delete group
        </button>
      </div>
      <div>
        <button onClick={() => setAddUserModalOpen(true)}>Add user</button>
      </div>
      <div>
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Balance</th>
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr key={user.id}>
                <td>{user.name}</td>
                <td>{user.balance}</td>
                <td>
                  <button
                    onClick={() => {
                      setSelectedUser(user);
                      setPatchUserModalOpen(true);
                    }}
                  >
                    Edit
                  </button>
                </td>
                <td>
                  <button
                    onClick={() => {
                      setSelectedUser(user);
                      setDeleteUserModalOpen(true);
                    }}
                  >
                    &times;
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
}
