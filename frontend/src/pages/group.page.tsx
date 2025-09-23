import { skipToken, useQuery } from '@tanstack/react-query';

import type { AxiosError } from 'axios';
import { useParams } from 'react-router';
import { getGroupById } from 'src/services/group.service';
import { getMembersByGroupId } from 'src/services/member.service';
import type { Group, Member } from 'src/types/common.type';

export function Group() {
  const { groupId = '' } = useParams();

  // group info
  const { data: group = null, isFetching: isFetchingGroup } = useQuery<
    Group,
    AxiosError
  >({
    queryKey: ['group', groupId],
    queryFn: () => getGroupById(Number(groupId)),
  });

  // member info
  const { data: members = [], isFetching: isFetchingMembers } = useQuery<
    Member[],
    AxiosError
  >({
    queryKey: ['members', groupId],
    queryFn: group ? () => getMembersByGroupId(group.id) : skipToken,
  });

  const isLoading = isFetchingGroup || isFetchingMembers;

  return !group || isLoading ? (
    <div>Loading...</div>
  ) : (
    <>
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
