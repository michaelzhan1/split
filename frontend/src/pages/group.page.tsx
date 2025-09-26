import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router';

import { skipToken, useMutation, useQuery } from '@tanstack/react-query';
import type { AxiosError } from 'axios';

import { AddMemberModal } from 'src/components/add-member-modal.component';
import { PatchGroupModal } from 'src/components/patch-group-modal.component';
import { getGroupById, patchGroup } from 'src/services/group.service';
import {
  addMembertoGroup,
  getMembersByGroupId,
} from 'src/services/member.service';
import type { Group, Member } from 'src/types/common.type';

export function Group() {
  const { groupId = '' } = useParams();
  const navigate = useNavigate();
  const [addMemberModalOpen, setAddMemberModalOpen] = useState<boolean>(false);
  const [patchGroupModalOpen, setPatchGroupModalOpen] =
    useState<boolean>(false);

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

  // member info
  const {
    data: members = [],
    isFetching: isFetchingMembers,
    refetch: refetchMembers,
    error: membersError,
  } = useQuery<Member[], AxiosError>({
    queryKey: ['members', groupId],
    queryFn: group ? () => getMembersByGroupId(group.id) : skipToken,
  });

  useEffect(() => {
    if (membersError) {
      console.error('Error fetching members:', membersError);
      alert('Failed to fetch members. Please try again.');
    }
  }, [membersError]);

  // add a member
  const { mutate: addMemberMutate, isPending: isPendingAddMember } =
    useMutation<{ id: number }, AxiosError, { name: string }>({
      mutationFn: (variables: { name: string }) => {
        return addMembertoGroup(Number(groupId), variables.name);
      },
    });
  const onAddMember = (name: string) =>
    addMemberMutate(
      { name },
      {
        onSuccess: () => {
          refetchMembers();
          setAddMemberModalOpen(false);
        },
        onError: (error) => {
          console.error('Error adding member:', error);
          alert('Failed to add member. Please try again');
        },
      },
    );

  const isLoading =
    isFetchingGroup ||
    isPendingPatchGroup ||
    isFetchingMembers ||
    isPendingAddMember;

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
      <AddMemberModal
        isOpen={addMemberModalOpen}
        onClose={() => setAddMemberModalOpen(false)}
        onSubmit={(name: string) => onAddMember(name)}
      />
      <div>
        <h1>Group: {group.name}</h1>
        <button onClick={() => setPatchGroupModalOpen(true)}>
          Edit group name
        </button>
      </div>
      <div>
        <button onClick={() => setAddMemberModalOpen(true)}>Add member</button>
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
            {members.map((member) => (
              <tr key={member.id}>
                <td>{member.name}</td>
                <td>{member.balance}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
}
