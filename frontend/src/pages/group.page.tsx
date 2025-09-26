import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router';

import { skipToken, useMutation, useQuery } from '@tanstack/react-query';
import type { AxiosError } from 'axios';

import { AddMemberModal } from 'src/components/add-member-modal.component';
import { getGroupById } from 'src/services/group.service';
import {
  addMembertoGroup,
  getMembersByGroupId,
} from 'src/services/member.service';
import type { Group, Member } from 'src/types/common.type';

export function Group() {
  const { groupId = '' } = useParams();
  const navigate = useNavigate();
  const [addMemberModalOpen, setAddMemberModalOpen] = useState<boolean>(false);

  // group info
  const {
    data: group = null,
    isFetching: isFetchingGroup,
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
      navigate('/');
    }
  }, [membersError, navigate]);

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

  const isLoading = isFetchingGroup || isFetchingMembers || isPendingAddMember;

  return !group || isLoading ? (
    <div>Loading...</div>
  ) : (
    <>
      <AddMemberModal
        isOpen={addMemberModalOpen}
        onClose={() => setAddMemberModalOpen(false)}
        onSubmit={(name: string) => onAddMember(name)}
      />
      <div>
        <h1>Group: {group.name}</h1>
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
