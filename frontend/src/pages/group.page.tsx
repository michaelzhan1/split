import { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router';

import { skipToken, useQuery } from '@tanstack/react-query';
import type { AxiosError } from 'axios';

import { Modal } from 'src/components/modal.component';
import { getGroupById } from 'src/services/group.service';
import { getMembersByGroupId } from 'src/services/member.service';
import type { Group, Member } from 'src/types/common.type';

export function Group() {
  const { groupId = '' } = useParams();
  const navigate = useNavigate();

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
    error: membersError,
  } = useQuery<Member[], AxiosError>({
    queryKey: ['members', groupId],
    queryFn: group ? () => getMembersByGroupId(group.id) : skipToken,
  });

  useEffect(() => {
    console.log(membersError);
    if (membersError) {
      console.error('Error fetching members:', membersError);
      alert('Failed to fetch members. Please try again.');
      navigate('/');
    }
  }, [membersError, navigate]);

  const isLoading = isFetchingGroup || isFetchingMembers;

  return !group || isLoading ? (
    <div>Loading...</div>
  ) : (
    <>
      <Modal
        isOpen={true}
        onClose={() => {}}
        title='Group Details'
        onSubmit={() => {}}
      >
        Test
      </Modal>
      <div>
        <h1>Group: {group.name}</h1>
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
